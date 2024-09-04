(function() {
    'use strict';
    const GOTO = 0;
    const GOBACK = 1;
    const SETVALUE = 2;
    const SETREPEATVALUE = 3;
    const UNSETVALUE = 4;
    function execute(root, steps) {
        let count = steps.length;
        let stack = [];
        let currRoot = root;
        for (let i = 0; i < count; i++) {
            let step = steps[i];
            // GOTO path
            if (step.type == GOTO) {
                let newRoot = currRoot;
                let path = "";
                for (const entry of step.path) {
                    path += "/" + entry;
                    if (newRoot[entry] == null) {
                        throw 'Invalid path ' + path;
                    } else  {
                        newRoot = newRoot[entry];
                    }
                }
                stack.push(currRoot);
                currRoot = newRoot;
            } else if (step.type == GOBACK) {
                currRoot = stack.pop();
            } else if (step.type == SETVALUE) {
                currRoot[step.key] = step.value; 
            } else if (step.type == SETREPEATVALUE) {
                for (const key of step.keys) {
                    currRoot[key] = step.value;
                }
            }
        }
    }

    function isParent(oldString, newString) {
        // The strings cannot be the same
        // and newString is not a substring
        if (oldString.length >= newString.length) {
            return false;
        }
        // Has to have a separate at the start.
        if (newString[oldString.length] !== "\x03") {
            return false;
        }

        for (let i = 0; i < oldString.length; i++) {
            if (oldString[i] != newString[i]) {
                return false;
            }
        }
        return true;
    }
    const tableRoot = Symbol('TableRoot');

    function buildSteps(stepTree, steps, patchPath = []) {
        let storedPatchPath = patchPath; 
        if (stepTree[tableRoot] === true) {
            steps.push({
                type: GOTO,
                path: Array.from(patchPath)
            });
            patchPath = [];
        }

        let setPairs = []; 
        let substeps = [];
        for (const pair of Object.entries(stepTree)) {
            const [key,value] = pair;
            if (key === tableRoot) {
               continue;
            }

            if (value == null) {
                substeps.push({
                    type: UNSETVALUE,
                    key: key
                });
            } else if (value.constructor === Object) {
                patchPath.push(key);
                buildSteps(value, steps, patchPath);
                patchPath.pop();
            } else {
                setPairs.push(pair);
            }
        }
        let matches = {};
        for(let i = 0; i < setPairs.length; i++) {
            const aPair = setPairs[i];
            const value = aPair[1];
            if (typeof value === "boolean" || 
                typeof value === "string" ||
                typeof value === "number") {
                if (matches[value] == null) {
                    matches[value] = [];
                }
                matches[value].push(i);
            } else {
                substeps.push({
                    type: SETVALUE,
                    key: aPair[0],
                    value,
                });
            }
        }
        for (const valueGroups of Object.values(matches)) {
            let firstIndex = valueGroups[0];
            let firstPair = setPairs[firstIndex];
        
            if (valueGroups.length === 1) {
                let key = firstPair[0];
                let value = firstPair[1];
                substeps.push({
                    type: SETVALUE,
                    key,
                    value,
                });
            } else {
                let groupValue = firstPair[1];
                substeps.push({
                    type: SETREPEATVALUE,
                    keys: valueGroups.map((i) => setPairs[i][0]),
                    value: groupValue
                });
            }
        }
        // TODO: Optimize substeps
        steps.push(...substeps);
        if (stepTree[tableRoot] === true) {
            steps.push({
                type: GOBACK
            });
            patchPath = storedPatchPath;
        }
    }

    function convertToSteps(json) {
        const steps = [];
        const root = {};
        const rootValues = {};
        let index = 1;
        // Convert the list of patches to a tree.
        // Use special value to denote the root 
        // of a keyed table
        for (const [key, value] of json) {
            let tempRoot = root;
            let subroots = key.slice(0, -1);
            let subkey = key.slice(-1);
            for (const subroot of subroots) {
                if (tempRoot[subroot] == null) {
                    tempRoot[subroot] = {};
                } else if (tempRoot[subroot].constructor !== Object) {
                    throw `Cannot access ${key.join('/')} a part of the path is not a table.`;
                }
                tempRoot = tempRoot[subroot];
            }
            // Only do it for
            // subrooted tables
            if (subroots.length >= 0) {
                tempRoot[tableRoot] = true;
            }
            tempRoot[subkey] = value;
        }
        buildSteps(root, steps);
        return steps;
    }
    return {execute, convertToSteps};
})()
