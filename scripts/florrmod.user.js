// ==UserScript==
// @name         florrmod API
// @namespace    https://github.com/z46-dev/florr-legacy
// @version      0.1.1
// @description  Exposes an API for interacting with the game state.
// @match        https://static.florr.io/old*
// @run-at       document-start
// @inject-into  page
// @grant        none
// ==/UserScript==

(() => {
    "use strict";

    if (window.florrmod) {
        return;
    }

    const NativeWorker = window.Worker;
    let serverWorker = null;

    const workerBootstrap = (serverURL, wasmURL) => {
        const nativeInstantiate = WebAssembly.instantiate.bind(WebAssembly);
        let serverInstance = null,
            resolveServerReady,
            afterSimulationTick = () => {};

        const nativeSetInterval = self.setInterval.bind(self);
        const nativeSetTimeout = self.setTimeout.bind(self);
        const runTimerCallback = (callback, callbackArgs) => {
            const result = callback(...callbackArgs);
            afterSimulationTick();
            return result;
        };

        self.setInterval = (callback, delay, ...args) => nativeSetInterval((...callbackArgs) => {
            return runTimerCallback(callback, callbackArgs);
        }, delay, ...args);
        self.setTimeout = (callback, delay, ...args) => nativeSetTimeout((...callbackArgs) => {
            return runTimerCallback(callback, callbackArgs);
        }, delay, ...args);

        const serverReady = new Promise((resolve) => resolveServerReady = resolve);

        // Utility functions for encoding and decoding unsigned LEB128 integers.
        const decodeULEB = (bytes, offset) => {
            let value = 0,
                shift = 0,
                cursor = offset;

            while (true) {
                const byte = bytes[cursor++];
                value |= (byte & 0x7f) << shift;

                if ((byte & 0x80) === 0) {
                    return {
                        value: value >>> 0,
                        next: cursor
                    };
                }

                shift += 7;
            }
        };

        // Utility function for encoding unsigned LEB128 integers.
        const encodeULEB = (value) => {
            const encoded = [];
            let remaining = value >>> 0;

            do {
                let byte = remaining & 0x7f;
                remaining >>>= 7;

                if (remaining !== 0) {
                    byte |= 0x80;
                }

                encoded.push(byte);
            } while (remaining !== 0);

            return encoded;
        };

        // Adds ordinary function exports without changing any function body or index.
        const exposeFunctions = (source) => {
            const bytes = new Uint8Array(source);
            const parts = [bytes.subarray(0, 8)];
            const exportsToAdd = [
                ["__florrmod_setHealth", 18],
                ["__florrmod_setMaxHealth", 25],
                ["__florrmod_setX", 44],
                ["__florrmod_setY", 47],
                ["__florrmod_setPetal", 125],
                ["__florrmod_spawnMob", 144],
                ["__florrmod_addExperience", 215],
            ];

            const newExports = exportsToAdd.flatMap(([name, index]) => {
                const encodedName = new TextEncoder().encode(name);
                return [...encodeULEB(encodedName.length), ...encodedName, 0, ...encodeULEB(index)];
            });

            let cursor = 8,
                foundExportSection = false;

            while (cursor < bytes.length) {
                const sectionStart = cursor;
                const sectionID = bytes[cursor++];
                const size = decodeULEB(bytes, cursor);
                const payloadStart = size.next;
                const payloadEnd = payloadStart + size.value;

                if (sectionID === 7) {
                    const count = decodeULEB(bytes, payloadStart);
                    const payload = [
                        ...encodeULEB(count.value + exportsToAdd.length),
                        ...bytes.subarray(count.next, payloadEnd),
                        ...newExports,
                    ];

                    parts.push(new Uint8Array([sectionID, ...encodeULEB(payload.length), ...payload]));
                    foundExportSection = true;
                } else {
                    parts.push(bytes.subarray(sectionStart, payloadEnd));
                }

                cursor = payloadEnd;
            }

            if (!foundExportSection) {
                throw new Error("florrmod: sp_server.wasm has no export section");
            }

            const output = new Uint8Array(parts.reduce((length, part) => length + part.length, 0));
            let outputOffset = 0;

            for (const part of parts) {
                output.set(part, outputOffset);
                outputOffset += part.length;
            }

            return output;
        };

        WebAssembly.instantiate = async (moduleOrBytes, imports) => {
            const result = await nativeInstantiate(moduleOrBytes, imports);
            serverInstance = result.instance ?? result;
            resolveServerReady();
            return result;
        };

        const wasmRequest = new XMLHttpRequest();
        wasmRequest.open("GET", wasmURL, false);
        wasmRequest.responseType = "arraybuffer";
        wasmRequest.send(null);

        if (wasmRequest.status < 200 || wasmRequest.status >= 300) {
            throw new Error(`florrmod: failed to load ${wasmURL} (${wasmRequest.status})`);
        }

        if (wasmRequest.response.byteLength !== 477558) {
            throw new Error(`florrmod: unsupported sp_server.wasm build (${wasmRequest.response.byteLength} bytes; expected 477558)`);
        }

        self.Module = {
            wasmBinary: exposeFunctions(wasmRequest.response),
            locateFile: (path) => path.endsWith(".wasm") ? wasmURL : new URL(path, serverURL).href,
        };

        importScripts(serverURL);
        const gameOnMessage = self.onmessage;

        const views = () => {
            if (!serverInstance) {
                throw new Error("florrmod: server WASM is not initialized yet");
            }

            const buffer = serverInstance.exports.o.buffer;
            return {
                bytes: new Uint8Array(buffer),
                u32: new Uint32Array(buffer),
                f32: new Float32Array(buffer),
                f64: new Float64Array(buffer),
            };
        };

        const readCString = (bytes, address) => {
            if (!address || address >= bytes.length) {
                return "";
            }

            let end = address;
            while (end < bytes.length && bytes[end] !== 0) {
                end += 1;
            }

            return new TextDecoder().decode(bytes.subarray(address, end));
        };

        const rarityNames = ["Common", "Unusual", "Rare", "Epic", "Legendary", "Unique"];

        const readPetals = () => {
            const { bytes, u32, f32 } = views();
            return Array.from({ length: 39 }, (_, index) => {
                const address = 8545792 + index * 52;
                const rarityCode = u32[(address + 24) >>> 2];
                return {
                    id: u32[(address + 48) >>> 2],
                    name: readCString(bytes, u32[address >>> 2]),
                    description: readCString(bytes, u32[(address + 8) >>> 2]),
                    rarity: { id: rarityCode, name: rarityNames[rarityCode] ?? null },
                    reloadMilliseconds: f32[(address + 28) >>> 2],
                    health: null,
                    damage: null,
                    entityCount: u32[(address + 32) >>> 2],
                    radius: f32[(address + 36) >>> 2],
                    flags: u32[(address + 40) >>> 2],
                    auxiliaryValue: f32[(address + 44) >>> 2],
                    unknown16: u32[(address + 16) >>> 2],
                    unknown20: u32[(address + 20) >>> 2],
                };
            });
        };

        // Each tuple is [entity component offset, active count index, active-list pointer index].
        const mobAIComponents = [
            [[108, 31, 33]], [[108, 31, 33]], [[176, 335, 337]], [[144, 139, 141]],
            [[144, 139, 141]], [[196, 430, 432]], [[132, 106, 108]], [[132, 106, 108]],
            [[132, 106, 108]], [[132, 106, 108]], [[136, 117, 119]], [[152, 269, 271]],
            [[152, 269, 271]], [[180, 346, 348]], [[172, 324, 326]], [[172, 324, 326]],
            [[172, 324, 326]], [[172, 324, 326]], [[128, 95, 97]], [[144, 139, 141]],
            [[144, 139, 141]], [[112, 42, 44], [128, 95, 97], [180, 346, 348]],
        ];

        const readMobs = () => {
            const { bytes, u32, f32 } = views();
            const petals = new Map(readPetals().map((petal) => [petal.id, petal]));
            return Array.from({ length: 22 }, (_, index) => {
                const address = 8548484 + index * 116;
                const dropCount = u32[(address + 32) >>> 2];
                const id = u32[(address + 100) >>> 2];
                const drops = Array.from({ length: dropCount }, (_, dropIndex) => {
                    const petalID = u32[(address + 36 + dropIndex * 8) >>> 2];
                    return {
                        petalID,
                        petalName: petals.get(petalID)?.name ?? null,
                        probability: f32[(address + 40 + dropIndex * 8) >>> 2],
                    };
                });

                const rarityCode = u32[(address + 24) >>> 2];
                return {
                    id,
                    registryIndex: index,
                    name: readCString(bytes, u32[address >>> 2]),
                    description: readCString(bytes, u32[(address + 8) >>> 2]),
                    rarity: { id: rarityCode, name: rarityNames[rarityCode] ?? null },
                    health: null,
                    damage: null,
                    scale: f32[(address + 112) >>> 2],
                    sizeHealthOrScoreScale: f32[(address + 28) >>> 2],
                    drops,
                    enabled: Boolean(u32[(address + 108) >>> 2]),
                    unknown12: u32[(address + 12) >>> 2],
                    unknown16: u32[(address + 16) >>> 2],
                    unknown20: u32[(address + 20) >>> 2],
                    unknown104: u32[(address + 104) >>> 2],
                };
            });
        };

        const playerContext = () => {
            const { bytes, u32, f32, f64 } = views();
            const activeServer = u32[24936 >>> 2];
            const session = activeServer ? u32[(activeServer + 120) >>> 2] : 0;
            const entity = session ? u32[(session + 100) >>> 2] : 0;
            const pointers = {
                entity,
                flower: entity ? u32[(entity + 200) >>> 2] : 0,
                health: entity ? u32[(entity + 164) >>> 2] : 0,
                position: entity ? u32[(entity + 168) >>> 2] : 0,
            };

            if (Object.values(pointers).some((pointer) => !pointer)) {
                return null;
            }

            const health = f32[(pointers.health + 8) >>> 2];
            if (!(health > 0)) {
                return null;
            }

            const level = bytes[pointers.flower + 148];
            const slotCount = level < 15 ? 13 : level < 30 ? 14 : level < 45 ? 15 : 16;
            return {
                pointers,
                info: {
                    id: u32[(entity + 16) >>> 2],
                    alive: true,
                    level,
                    experience: f32[(pointers.flower + 16) >>> 2],
                    totalExperience: f32[(pointers.flower + 280) >>> 2],
                    health,
                    maxHealth: f32[(pointers.health + 12) >>> 2],
                    position: {
                        x: f64[(pointers.position + 88) >>> 3],
                        y: f64[(pointers.position + 96) >>> 3],
                    },
                    loadout: Array.from(bytes.subarray(pointers.flower + 20, pointers.flower + 20 + slotCount)),
                },
            };
        };

        const getPlayer = () => playerContext()?.info ?? null;

        const requirePlayer = () => {
            const player = playerContext();
            if (!player) {
                throw new Error("florrmod: the player is not alive");
            }

            return { me: player.info, pointers: player.pointers };
        };

        const finite = (value, name) => {
            if (!Number.isFinite(value)) {
                throw new TypeError(`florrmod: ${name} must be a finite number`);
            }

            return value;
        };

        const coordinates = (position) => {
            if (!position || typeof position !== "object") {
                throw new TypeError("florrmod: position must be an object with x and y numbers");
            }

            return {
                x: finite(position.x, "position.x"),
                y: finite(position.y, "position.y")
            };
        };

        const withScratch = (size, callback) => {
            const pointer = serverInstance.exports.q(size);
            try {
                return callback(pointer);
            } finally {
                serverInstance.exports.r(pointer);
            }
        };

        const setPosition = (component, position) => {
            const { x, y } = coordinates(position);
            withScratch(16, (scratch) => {
                const { f64 } = views();
                f64[scratch >>> 3] = x;
                f64[(scratch + 8) >>> 3] = y;
                serverInstance.exports.__florrmod_setX(component + 20, scratch);
                serverInstance.exports.__florrmod_setY(component + 16, scratch + 8);
            });
        };

        const setHealthField = (component, offset, value) => {
            withScratch(4, (scratch) => {
                views().f32[scratch >>> 2] = value;
                const setter = offset === 8 ? "__florrmod_setHealth" : "__florrmod_setMaxHealth";
                serverInstance.exports[setter](component + offset, scratch);
            });
        };

        const spawnedMobs = new Map();

        const mobContext = (id) => {
            const tracked = spawnedMobs.get(id);
            if (!tracked) {
                return null;
            }

            const { u32, f32, f64 } = views();
            if (u32[(tracked.entity + 16) >>> 2] !== id) {
                spawnedMobs.delete(id);
                return null;
            }

            const pointers = {
                health: u32[(tracked.entity + 164) >>> 2],
                position: u32[(tracked.entity + 168) >>> 2],
                collision: u32[(tracked.entity + 188) >>> 2],
            };
            if (Object.values(pointers).some((pointer) => !pointer)) {
                spawnedMobs.delete(id);
                return null;
            }

            const health = f32[(pointers.health + 8) >>> 2];
            return {
                tracked,
                pointers,
                info: {
                    id,
                    mobID: tracked.mobID,
                    registryIndex: tracked.registryIndex,
                    name: tracked.name,
                    alive: health > 0,
                    aiEnabled: tracked.aiEnabled,
                    health,
                    maxHealth: f32[(pointers.health + 12) >>> 2],
                    size: f32[(pointers.collision + 8) >>> 2],
                    _nativeSize: tracked.nativeSize,
                    _playerPosition: getPlayer()?.position ?? null,
                    position: {
                        x: f64[(pointers.position + 88) >>> 3],
                        y: f64[(pointers.position + 96) >>> 3],
                    },
                },
            };
        };

        const requireMob = (id) => {
            if (!Number.isInteger(id)) {
                throw new TypeError("florrmod: mob entity id must be an integer");
            }

            const mob = mobContext(id);
            if (!mob || !mob.info.alive) {
                throw new Error(`florrmod: spawned mob ${id} is no longer alive`);
            }

            return mob;
        };

        const setMobAI = (tracked, enabled) => {
            const { u32, f64 } = views();
            const world = u32[tracked.entity >>> 2];
            for (const component of tracked.aiComponents) {
                const countAddress = (world + component.countIndex * 4) >>> 2;
                const list = u32[(world + component.listIndex * 4) >>> 2];
                let count = u32[countAddress];

                if (enabled && !component.active) {
                    const capacity = u32[(world + (component.listIndex - 1) * 4) >>> 2];
                    if (count >= capacity) {
                        throw new Error("florrmod: the mob AI list filled while this component was disabled");
                    }

                    u32[(list + count * 4) >>> 2] = component.pointer;
                    u32[countAddress] = count + 1;
                    component.active = true;
                } else if (!enabled && component.active) {
                    for (let index = 0; index < count; index += 1) {
                        if (u32[(list + index * 4) >>> 2] !== component.pointer) {
                            continue;
                        }

                        count -= 1;
                        u32[(list + index * 4) >>> 2] = u32[(list + count * 4) >>> 2];
                        u32[countAddress] = count;
                        break;
                    }

                    component.active = false;
                }
            }

            tracked.aiEnabled = enabled;
            if (enabled) {
                tracked.lockedPosition = null;
            } else {
                const position = u32[(tracked.entity + 168) >>> 2];
                tracked.lockedPosition = {
                    x: f64[(position + 88) >>> 3],
                    y: f64[(position + 96) >>> 3],
                };
            }
        };

        afterSimulationTick = () => {
            for (const [id, tracked] of spawnedMobs) {
                if (tracked.aiEnabled || !tracked.lockedPosition) {
                    continue;
                }

                const mob = mobContext(id);
                if (mob?.info.alive) {
                    setPosition(mob.pointers.position, tracked.lockedPosition);
                    const { bytes, f64 } = views();
                    f64[(mob.pointers.position + 88) >>> 3] = tracked.lockedPosition.x;
                    f64[(mob.pointers.position + 96) >>> 3] = tracked.lockedPosition.y;
                    bytes[mob.pointers.position + 80] = 0;
                    bytes[mob.pointers.collision + 16] = 0;
                }
            }
        };

        const setMobProperties = (mob, properties) => {
            if (!properties || typeof properties !== "object" || Array.isArray(properties)) {
                throw new TypeError("florrmod: properties must be an object");
            }

            const supported = new Set(["aiEnabled", "health", "maxHealth", "position"]);
            for (const property of Object.keys(properties)) {
                if (!supported.has(property)) {
                    throw new RangeError(`florrmod: unsupported mob property ${property}`);
                }
            }

            let maxHealth = mob.info.maxHealth,
                position = null;
            if (properties.maxHealth !== undefined) {
                finite(properties.maxHealth, "maxHealth");
                if (properties.maxHealth <= 0) {
                    throw new RangeError("florrmod: maxHealth must be greater than zero");
                }

                maxHealth = properties.maxHealth;
                setHealthField(mob.pointers.health, 12, maxHealth);
            }

            if (properties.health !== undefined) {
                finite(properties.health, "health");
                if (properties.health <= 0) {
                    throw new RangeError("florrmod: health must be greater than zero; zero bypasses the game's death path");
                }

                setHealthField(mob.pointers.health, 8, Math.min(properties.health, maxHealth));
            } else if (mob.info.health > maxHealth) {
                setHealthField(mob.pointers.health, 8, maxHealth);
            }

            if (properties.position !== undefined) {
                position = coordinates(properties.position);
                setPosition(mob.pointers.position, position);
            }

            if (properties.aiEnabled !== undefined) {
                if (typeof properties.aiEnabled !== "boolean") {
                    throw new TypeError("florrmod: aiEnabled must be a boolean");
                }

                setMobAI(mob.tracked, properties.aiEnabled);
            }

            if (position && !mob.tracked.aiEnabled) {
                mob.tracked.lockedPosition = position;
            }
        };

        const spawnMob = (mobID, position, properties = {}) => {
            requirePlayer();
            const mobs = readMobs();
            const registryIndex = Number.isInteger(mobID) ? mobs.findIndex((mob) => mob.id === mobID && mob.enabled) : -1;
            if (registryIndex < 0) {
                throw new RangeError(`florrmod: unknown or disabled mob id ${mobID}`);
            }

            const { x, y } = coordinates(position);
            const entity = serverInstance.exports.__florrmod_spawnMob(24864, registryIndex + 1, x, y, -1);
            if (!entity) {
                throw new Error("florrmod: the game entity pool could not allocate a mob");
            }

            const { u32 } = views();
            const id = u32[(entity + 16) >>> 2];
            const aiComponents = mobAIComponents[registryIndex].map(([offset, countIndex, listIndex]) => ({
                pointer: u32[(entity + offset) >>> 2],
                countIndex,
                listIndex,
                active: true,
            })).filter(({ pointer }) => pointer);

            spawnedMobs.set(id, {
                id,
                entity,
                mobID,
                registryIndex,
                name: mobs[registryIndex].name,
                aiEnabled: true,
                aiComponents,
            });

            const mob = requireMob(id);
            mob.tracked.nativeSize = mob.info.size;
            setMobProperties(mob, properties);
            return mobContext(id).info;
        };

        const handlers = {
            getPetals: readPetals,
            getMobs: readMobs,
            getMe: getPlayer,
            setMyPetal({ petalID, slotID }) {
                const { me, pointers } = requirePlayer();
                if (!Number.isInteger(slotID) || slotID < 0 || slotID >= me.loadout.length) {
                    throw new RangeError(`florrmod: slotID must be between 0 and ${me.loadout.length - 1}`);
                }

                if (!readPetals().some((petal) => petal.id === petalID)) {
                    throw new RangeError(`florrmod: unknown petal id ${petalID}`);
                }

                serverInstance.exports.__florrmod_setPetal(pointers.flower, slotID, petalID);
                return getPlayer();
            },
            setMyPetals({ petalIDs }) {
                const { me, pointers } = requirePlayer();
                if (!Array.isArray(petalIDs) || petalIDs.length !== me.loadout.length) {
                    throw new RangeError(`florrmod: petalIDs must contain exactly ${me.loadout.length} ids`);
                }

                const knownPetals = new Set(readPetals().map((petal) => petal.id));
                petalIDs.forEach((petalID, slotID) => {
                    if (!knownPetals.has(petalID)) {
                        throw new RangeError(`florrmod: unknown petal id ${petalID} at slot ${slotID}`);
                    }
                });

                petalIDs.forEach((petalID, slotID) => {
                    serverInstance.exports.__florrmod_setPetal(pointers.flower, slotID, petalID);
                });

                return getPlayer();
            },
            setMe({ properties }) {
                const { me, pointers } = requirePlayer();
                if (!properties || typeof properties !== "object" || Array.isArray(properties)) {
                    throw new TypeError("florrmod: properties must be an object");
                }

                const supported = new Set(["health", "maxHealth", "position"]);
                for (const property of Object.keys(properties)) {
                    if (!supported.has(property)) {
                        throw new RangeError(`florrmod: unsupported player property ${property}`);
                    }
                }

                let maxHealth = me.maxHealth;
                if (properties.maxHealth !== undefined) {
                    finite(properties.maxHealth, "maxHealth");
                    if (properties.maxHealth <= 0) {
                        throw new RangeError("florrmod: maxHealth must be greater than zero");
                    }

                    maxHealth = properties.maxHealth;
                    setHealthField(pointers.health, 12, maxHealth);
                }

                if (properties.health !== undefined) {
                    finite(properties.health, "health");
                    if (properties.health <= 0) {
                        throw new RangeError("florrmod: health must be greater than zero; zero bypasses the game's death path");
                    }

                    setHealthField(pointers.health, 8, Math.min(properties.health, maxHealth));
                } else if (me.health > maxHealth) {
                    setHealthField(pointers.health, 8, maxHealth);
                }

                if (properties.position !== undefined) {
                    setPosition(pointers.position, properties.position);
                }

                return getPlayer();
            },
            addMyExperience({ amount }) {
                const { pointers } = requirePlayer();
                finite(amount, "amount");
                if (amount < 0) {
                    throw new RangeError("florrmod: amount must not be negative");
                }

                serverInstance.exports.__florrmod_addExperience(pointers.flower, amount);
                return getPlayer();
            },
            getMob({ id }) {
                return mobContext(id)?.info ?? null;
            },
            getSpawnedMobs() {
                return Array.from(spawnedMobs.keys(), (id) => mobContext(id)?.info).filter(Boolean);
            },
            setMob({ id, properties }) {
                const mob = requireMob(id);
                setMobProperties(mob, properties);
                return mobContext(id).info;
            },
            spawnMob({ mobID, position, properties }) {
                return spawnMob(mobID, position, properties);
            },
            spawnDummy({ mobID, position, properties }) {
                return spawnMob(mobID, position, {
                    maxHealth: 1000000000,
                    health: 1000000000,
                    ...properties,
                    aiEnabled: false,
                });
            },
        };

        self.onmessage = (event) => {
            if (!event.data?.__florrmod) {
                gameOnMessage(event);
                return;
            }

            const port = event.ports[0];
            Promise.resolve().then(() => serverReady).then(() => handlers[event.data.method](event.data.args ?? {})).then((value) => port.postMessage({ value })).catch((error) => port.postMessage({
                error: {
                    name: error.name,
                    message: error.message,
                    stack: error.stack
                }
            }));
        };
    };

    const bootstrapSource = `(${workerBootstrap.toString()})(SERVER_URL, WASM_URL);`;

    const hitboxMobs = new Map();
    let showHitboxes = false;
    let renderFrame = 0;

    const matchHitboxMob = (context, transform) => {
        if (!showHitboxes || hitboxMobs.size === 0) {
            return null;
        }

        const zoom = Math.hypot(transform[0], transform[1]);
        if (!(zoom > 0.05 && zoom < 5)) {
            return null;
        }

        for (const target of hitboxMobs.values()) {
            if (target.drawFrame === renderFrame || !target.playerPosition) {
                continue;
            }

            const centerX = context.canvas.width / 2 + (target.position.x - target.playerPosition.x) * zoom;
            const centerY = context.canvas.height / 2 + (target.position.y - target.playerPosition.y) * zoom;
            const tolerance = Math.max(4, target.nativeSize * zoom * 0.35);
            if (Math.hypot(transform[4] - centerX, transform[5] - centerY) <= tolerance) {
                target.drawFrame = renderFrame;
                return { x: transform[4], y: transform[5], radius: target.size * zoom };
            }
        }

        return null;
    };

    const drawHitbox = (context, hitbox) => {
        const transform = context.getTransform();
        context.setTransform(1, 0, 0, 1, 0, 0);
        context.save();
        context.lineWidth = 2;
        context.strokeStyle = "rgba(255, 64, 64, 0.95)";
        context.fillStyle = "rgba(255, 64, 64, 0.08)";
        context.beginPath();
        context.arc(hitbox.x, hitbox.y, hitbox.radius, 0, Math.PI * 2);
        context.fill();
        context.stroke();
        context.restore();
        context.setTransform(transform);
    };

    const installHitboxOverlay = (prototype) => {
        if (!prototype) {
            return;
        }

        const nativeSetTransform = prototype.setTransform;
        prototype.setTransform = function(...transform) {
            if (transform.length !== 6 || !transform.every(Number.isFinite)) {
                return nativeSetTransform.apply(this, transform);
            }

            const result = nativeSetTransform.apply(this, transform);
            const hitbox = matchHitboxMob(this, transform);
            if (hitbox) {
                drawHitbox(this, hitbox);
            }

            return result;
        };
    };

    installHitboxOverlay(window.CanvasRenderingContext2D?.prototype);
    installHitboxOverlay(window.OffscreenCanvasRenderingContext2D?.prototype);
    const advanceRenderFrame = () => {
        renderFrame += 1;
        requestAnimationFrame(advanceRenderFrame);
    };
    requestAnimationFrame(advanceRenderFrame);

    function FlorrmodWorker(url, options) {
        if (!new.target) {
            throw new TypeError("Failed to construct 'Worker': Please use the 'new' operator");
        }

        const resolvedURL = new URL(String(url), location.href);
        if (!resolvedURL.pathname.endsWith("/sp_server.js")) {
            return Reflect.construct(NativeWorker, [url, options], new.target);
        }

        const wasmURL = new URL("sp_server.wasm", resolvedURL).href;
        const source = bootstrapSource.replace("SERVER_URL", JSON.stringify(resolvedURL.href)).replace("WASM_URL", JSON.stringify(wasmURL));
        const blobURL = URL.createObjectURL(new Blob([source], { type: "text/javascript" }));
        const worker = new NativeWorker(blobURL, options);
        URL.revokeObjectURL(blobURL);
        serverWorker = worker;
        return worker;
    }

    Object.setPrototypeOf(FlorrmodWorker, NativeWorker);
    FlorrmodWorker.prototype = NativeWorker.prototype;
    window.Worker = FlorrmodWorker;

    const callServer = (method, args = {}) => new Promise((resolve, reject) => {
        if (!serverWorker) {
            reject(new Error("florrmod: the single-player server worker has not been created yet"));
            return;
        }

        const channel = new MessageChannel();
        const timeout = setTimeout(() => {
            channel.port1.close();
            reject(new Error(`florrmod: ${method} timed out waiting for the server worker`));
        }, 15000);

        channel.port1.onmessage = ({ data }) => {
            clearTimeout(timeout);
            channel.port1.close();

            if (data.error) {
                const error = new Error(data.error.message);
                error.name = data.error.name;
                error.stack = data.error.stack;
                reject(error);
            } else {
                resolve(data.value);
            }
        };

        serverWorker.postMessage({ __florrmod: true, method, args }, [channel.port2]);
    });

    const synchronizeMobHitbox = (info) => {
        if (!info) {
            return info;
        }

        const { _nativeSize: nativeSize, _playerPosition: playerPosition, ...publicInfo } = info;
        if (nativeSize > 0) {
            hitboxMobs.set(publicInfo.id, {
                id: publicInfo.id,
                registryIndex: publicInfo.registryIndex,
                nativeSize,
                size: publicInfo.size,
                position: publicInfo.position,
                playerPosition,
                drawFrame: -1,
            });
        } else {
            hitboxMobs.delete(publicInfo.id);
        }

        return publicInfo;
    };

    const callMobServer = async (method, args) => synchronizeMobHitbox(await callServer(method, args));

    Object.defineProperty(window, "florrmod", {
        configurable: false,
        enumerable: true,
        writable: false,
        value: Object.freeze({
            version: "0.1.1",
            getPetals: () => callServer("getPetals"),
            getMobs: () => callServer("getMobs"),
            getMe: () => callServer("getMe"),
            setMyPetal: (petalID, slotID) => callServer("setMyPetal", { petalID, slotID }),
            setMyPetals: (petalIDs) => callServer("setMyPetals", { petalIDs }),
            setMe: (properties) => callServer("setMe", { properties }),
            setMyHealth: (health) => callServer("setMe", { properties: { health } }),
            setMyMaxHealth: (maxHealth) => callServer("setMe", { properties: { maxHealth } }),
            setMyPosition: (x, y) => callServer("setMe", { properties: { position: { x, y } } }),
            addMyExperience: (amount) => callServer("addMyExperience", { amount }),
            getMob: (id) => callMobServer("getMob", { id }),
            getSpawnedMobs: async () => (await callServer("getSpawnedMobs")).map(synchronizeMobHitbox),
            setMob: (id, properties) => callMobServer("setMob", { id, properties }),
            setMobAI: (id, aiEnabled) => callMobServer("setMob", { id, properties: { aiEnabled } }),
            setMobHealth: (id, health) => callMobServer("setMob", { id, properties: { health } }),
            setMobMaxHealth: (id, maxHealth) => callMobServer("setMob", { id, properties: { maxHealth } }),
            setMobPosition: (id, x, y) => callMobServer("setMob", { id, properties: { position: { x, y } } }),
            ShowHitboxes: (enabled) => {
                showHitboxes = Boolean(enabled);
                return showHitboxes;
            },
            showHitboxes: (enabled) => {
                showHitboxes = Boolean(enabled);
                return showHitboxes;
            },
            spawnMob: (mobID, position, properties) => callMobServer("spawnMob", { mobID, position, properties }),
            spawnDummy: (mobID, position, properties) => callMobServer("spawnDummy", { mobID, position, properties }),
        }),
    });
})();
