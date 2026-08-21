package game

import (
	"github.com/z46-dev/arctic"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/configs"
)

func NewClient(g *Game, peer *arctic.ServerClient) (client *Client) {
	client = &Client{
		ID:     g.nextClientID,
		Peer:   peer,
		Game:   g,
		Camera: nil,
		Body:   nil,
	}

	client.Camera = NewCamera(client)
	g.nextClientID++
	g.clients.Add(client)
	return
}

func (c *Client) Update() {
	if c.Camera != nil {
		c.Camera.See()
	}
}

func (c *Client) OnMessage(data []byte) (err error) {
	var (
		reader  *shared.Reader = shared.NewReader(data)
		msgType uint8          = reader.GetU8()
	)

	switch msgType {
	case shared.GameMsgTypeSpawn:
		if c.Body != nil && !c.Body.isDead {
			return
		}

		c.Body = NewFlower(c.Game)
		c.Body.Name = reader.GetStringUTF8()
		if len(c.Body.Name) > shared.NameLengthLimit {
			c.Body.Name = c.Body.Name[:shared.NameLengthLimit]
		}

		c.Game.L.Infof("client %d spawned as %s\n", c.ID, c.Body.Name)
	case shared.GameMsgTypeInput:
		if c.Body == nil || c.Body.isDead {
			return
		}

		var angle, magnitude float64 = float64(reader.GetF32()), float64(reader.GetF32())
		c.Body.Movement.Zero().DirectionMagnitudeInto(angle, magnitude)
		c.Body.Facing = angle

		var action uint8 = reader.GetU8()
		switch action {
		case 0:
			c.Body.AttackMode = configs.AttackMode_Neutral
		case 1:
			c.Body.AttackMode = configs.AttackMode_Attack
		case 2:
			c.Body.AttackMode = configs.AttackMode_Defend
		}
	}

	return
}

func (c *Client) Remove() {}
