(function() {

    if (typeof window !== "undefined" && !window.decodeText) {
        window.decodeText = function(arr) {
            return new TextDecoder().decode(arr);
        }
    }

    function getString(view, off, end = -1) {
        if (end == -1) {
            end = off;
            for (; end < view.byteLength; end++){
                if (view.getInt8(end) == 0) {
                    break;
                }
            }
        }
        const textView = new Uint8Array(view.buffer);
        const str = textView.subarray(off, end);
        try {
            let result = decodeText(str);
            return result;
        } catch (e) {}

        return new Uint8Array(str);
    }
    
    function getVarName(view, chunkOffset) {
        let off = view.getInt32(chunkOffset, true);
        return getString(view, off);
    }
    
    function getType(view, chunkOffset) {
        return view.getInt32(chunkOffset + 0x4, true) & 7;
    }
    
    function getValueOffset(view, chunkOffset) {
        return view.getInt32(chunkOffset + 0x4, true) >> 4;
    }
    
    function readVector(view, vectorOffset) {
        return Array(4).fill(0).map((_, i) => {
            let off = vectorOffset + 0x4 * i
            return view.getFloat32(off, true);
        });
    }
    
    function parseDict(root, view, count, metadata) {
        for(let i = 0; i < count; i++) {
            parseChunk(root, view, metadata)
        }
    }
    
    const TDICT = 0
    const TFLOAT32 = 1
    const TVECTOR = 2
    const TSTRING = 3
    const TBOOLEAN = 4
    const TUNKNOWN = -1
    function parseChunk(root, view, metadata) {
        metadata["count"] += 1
        let chunkOffset = metadata["offset"]
        metadata["offset"] += 0x10
        const varName = getVarName(view, chunkOffset)
        const dataType = getType(view, chunkOffset)
        const valueOffset = getValueOffset(view, chunkOffset)
        // Empty BigInt used as data type to signify that
        if (view.getUint8(chunkOffset + 0x4) & 8 == 0) {
            root[varName] = BigInt(dataType)
            return
        }
        
        if (dataType == TDICT) {
            let newRoot = {}
            let count = view.getInt32(chunkOffset + 0x8, true)
            parseDict(newRoot, view, count, metadata)
            root[varName] = newRoot
        } else if (dataType == TFLOAT32) {
            root[varName] = view.getFloat32(chunkOffset + 0x8, true)
        } else if (dataType == TVECTOR) {
            root[varName] = readVector(view, valueOffset);
        } else if (dataType == TSTRING) {
            let stringLength = view.getInt32(chunkOffset + 0x8, true)
            root[varName] = getString(view, valueOffset, valueOffset + stringLength - 1)
        } else if (dataType == TBOOLEAN) {
            root[varName] = view.getUint8(chunkOffset + 0x8) > 0
        } else {
            throw ("Unknown data type " +  dataType)
        }
    }

    function convertToJson(buffer) {
        const view = new DataView(buffer);
        let numOfData = view.getInt32(0xC, true);
        const metadata = {
            "offset": 0x10,
            "count": 0
        };
        const root = {};
        while (metadata["count"] < numOfData) {
            parseChunk(root, view, metadata)
        }
        return root;
    }
    return convertToJson;
})();
