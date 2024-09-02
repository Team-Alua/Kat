(function() {

    'use strict';
    if (typeof window !== "undefined" && !window.encodeText) {
        window.encodeText = function(str) {
            return new TextEncoder().encode(str);
        }
    }


    function writeByteArray(view, offset, arr) {
        let iarr = new Uint8Array(view.buffer);
        copyBuffer(iarr, offset, arr)
    }

    function writeVector(view, offset, vector) {
        for (let i = 0; i < 4; i++) {
            view.setFloat32(offset + i * 4, vector[i], true)
        }
    }

    function writeUtf8String(view, offset, str) {
        let buffer = encodeText(str);
        let iarr = new Uint8Array(view.buffer);
        copyBuffer(iarr, offset, buffer)
    }

    function writeVarName(view, varName, writePointers) {
        let chunkOffset = writePointers["data"]
        let stringOffset = writePointers["strings"]
        view.setInt32(chunkOffset, stringOffset, true)
        writeUtf8String(view, stringOffset, varName)
        writePointers["strings"] += varName.length + 1
    }

    function writeHash(view, varName, writePointers) {
        let hash = fnv1a32(varName)
        let chunkOffset = writePointers["data"]
        view.setInt32(chunkOffset + 0xC, hash, true)
    }

    function serializeJson(view, rootJson, writePointers) {
        for (const [key, value] of Object.entries(rootJson)) {
            serializeChunk(view, key, value, writePointers)    
        }
    }

    function serializeChunk(view, varName, value, writePointers) {
        writeVarName(view, varName, writePointers)
        writeHash(view, varName, writePointers)
        let chunkOffset = writePointers["data"]
        writePointers["data"] += 0x10

        if (typeof value === "boolean") {
            view.setInt32(chunkOffset + 0x4, 0xC, true)
            view.setInt32(chunkOffset + 0x8, value ? 1 : 0, true)
        } else if (value instanceof Uint8Array) {
            let stringOffset = writePointers["strings"]
            writeByteArray(view, stringOffset, value)
            writePointers["strings"] += value.byteLength + 1
            let typeValue = (stringOffset << 4) + 0xB
            view.setUint32(chunkOffset + 4, typeValue, true)
            view.setInt32(chunkOffset + 8, value.byteLength + 1, true)
        } else if (typeof value === "string") {
            let stringOffset = writePointers["strings"]
            writeUtf8String(view, stringOffset, value)
            writePointers["strings"] += value.length + 1
            let typeValue = (stringOffset << 4) + 0xB
            view.setUint32(chunkOffset + 4, typeValue, true)
            view.setInt32(chunkOffset + 8, value.length + 1, true)
            // Also string
        } else if (Array.isArray(value)) {
            // Vector
            let vectorOffset = writePointers["vectors"]
            writeVector(view, vectorOffset, value)
            let typeValue = (vectorOffset << 4) + 0xA
            writePointers["vectors"] += 0x10
            view.setUint32(chunkOffset + 4, typeValue, true)
            view.setInt32(chunkOffset + 8, 0x10, true)
        } else if (typeof value === "number") {
            // float32
            view.setInt32(chunkOffset + 4, 0x9, true)
            view.setFloat32(chunkOffset + 8, value, true)
        } else if (typeof value === "object") {
            view.setInt32(chunkOffset + 4, 0x8, true)
            view.setInt32(chunkOffset + 8, Object.keys(value).length, true)
            serializeJson(view, value, writePointers)
        } else if (typeof value === "bigint") {
            // null value
            view.seInt32(chunkOffset + 4, Number(value), true)
        } else {
            throw ("Unknown data type " +  dataType)
        }
    }

    function analyzeJson(json, results) {
        for (let [key, value] of Object.entries(json)) {
            const unique = results["unique"]
            results["stringTotalSize"] += key.length + 1
            results["entries"] += 1
            if (Array.isArray(value)) {
                if (value.length == 4) {
                    results.vectors += 1
                }
            } else if (typeof value === "string") {
                results.stringTotalSize += value.length + 1
            } else if (value instanceof Uint8Array) {
                results.stringTotalSize += value.byteLength + 1
            } else if (typeof value === "object") {
                analyzeJson(value, results)
            }

        }
    }

    function fromJsonToSaveBin(json) {
        const magic = new Uint8Array([0x67, 0x67, 0x64, 0x4C, 0x89, 0x06, 0x33, 0x01])
        const info = {
            "entries": 0,
            "vectors": 0,
            "stringTotalSize": 0,
        }
        analyzeJson(json, info)
        let totalFileSize = (info["entries"] + info["vectors"] + 1) * 0x10 
        totalFileSize += info["stringTotalSize"]
        const buffer = new ArrayBuffer(totalFileSize);
        const view = new DataView(buffer);
        writeByteArray(view, 0x0, magic);
        view.setInt32(0x8, totalFileSize, true);
        view.setInt32(0xC, info["entries"], true);
        const sectionOffsets = {
            "data": 0x10,
            "vectors": 0x10 * (info["entries"] + 1),
            "strings": 0x10 * (info["entries"] + info["vectors"] + 1)
        };
        for (const [key, value] of Object.entries(json)) {
            serializeChunk(view, key, value, sectionOffsets)
        }
        return buffer;
    }
    return fromJsonToSaveBin;
})();
