(function() {
    'use strict';
    const GOTO = 0;
    const GOBACK = 1;
    const SETVALUE = 2;
    const SETREPEATVALUE = 3;
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
                    if (typeof root[entry] == "object") {
                        newRoot = root[entry];
                    } else {
                        throw 'Invalid path ' + path;
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

    function mergeToSteps(json) {
        const steps = [];
        const roots = {};
        const rootValues = {};
        let index = 1;
        // Find all roots
        for (const [key, value] of json) {
            let newKey = Array.from(key);
            let rootKey = newKey.pop();
            let root = newKey.join("\x03");
            if (roots[root] == null) {
                rootValues[root] = [];
                roots[root] = index++;
            }
            rootValues[root].push([rootKey,value])
        }
        // This is the order we are going to traverse the tree
        let sortedRoots = Object.keys(roots).sort();

        // TODO: Figure out algorithm to
        // track common parent between children
        for (let i = 0; i < sortedRoots.length; i++) {
            // look for matching keys
            let root = sortedRoots[i];
            let sameKeys = {};
            let values = rootValues[root];
            for (const valuePair of values) {
                const [rootKey, value] = valuePair;
                if (sameKeys[value] == null) {
                    sameKeys[value] = [];
                }
                sameKeys[value].push(valuePair);
            }
            let isEmptyRoot = root === "";
            if (!isEmptyRoot) {
                steps.push({
                    step: GOTO,
                    path: root.split("\x03")
                })
            }
            for (const [sharedKey, values] of Object.entries(sameKeys)) {
                if (values.length > 1) {
                    const keys = values.map(e => e[0]);
                    steps.push({
                        type: SETREPEATVALUE,
                        keys,
                        value: values[0][1]
                    });

                } else {
                    const [key, value] = values[0];
                    steps.push({
                        type: SETVALUE,
                        key,
                        value
                    });
                }
            }
            if (!isEmptyRoot) {
                steps.push({
                    step: GOBACK
                })
            }
        }
        return steps;
    }
    return {execute, mergeToSteps};
})()
