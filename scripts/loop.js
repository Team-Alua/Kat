function walker(root) {
    let files = [];
	let dirs = [root];
	while (dirs.length) {
		let root = dirs.shift()
		let data = fs.readdir(root)
		for (const {Name, Dir, Size} of data) {
			let p = root + Name
			if (Dir) {
				p += "/"
				dirs.push(p)
			} else {
                files.push({root: root, path: p, name: Name, size: Size})
            }
		}
	}
    return files;
}

function checkFiles(files) {
    let imgPath = "";
    let binPath = "";
    let parentPath = "";
    let imgName = "";
    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        if (file.name.endsWith(".bin")) {
            if (file.size != 96) {
                throw (`File size of ${file.name} should be 96 but is ${file.size}`);
            }
            binPath = file.path;
        } else {
            const size = file.size;
            const saveBlocks = 1 << 15;
            if (size%saveBlocks != 0) {
                throw (`Unexpected file size ${file.size} for ${file.name}`);
                // Not valid
            }
            const minImageSize = 96 * saveBlocks;
            const maxImageSize = (1 << 15) * saveBlocks;
            if (size > maxImageSize || size < minImageSize) {
                throw (`Unexpected file size ${file.size} for ${file.name}`);
            }
            imgName = file.name;
            imgPath = file.path;
            parentPath = file.root;
        }
    }
    if (imgPath + ".bin" != binPath) {
        throw (`${binName} does not go with ${imgName}`);
    }
    return {root: parentPath, name: imgName};
}

discord.sendMessage("Send your gravity rush save zip.\nIt should not be more than 4MBs.")
let msg = discord.getMessage()
download("/local/download.zip", msg.Attachments[0].URL)
fs.mount("/local/download.zip", "/mnt1", {"MountType": "zipfs", "ReadOnly": true})

fs.copyDir("/mnt1", "/tmp")

fs.unmount("/mnt1")

fs.remove("/local/download.zip")

let files = walker("/tmp/");
if (files.length != 2) {
    discord.sendMessage("There should only be a data000X and data000X.bin inside the zip.")
    exit()
}
const {root, name} = checkFiles(files);
console.log(msg)
const PS4RelRoot = msg.Message.Author.ID;
const userPS4Root = `/ps4/${PS4RelRoot}/`;
fs.mkdir(userPS4Root)
fs.mkdir(userPS4Root + "extract/")
fs.copyDir(root, userPS4Root)


const PS4Root = `/hostapp/${PS4RelRoot}`;

const cmds = [{
    "RequestType": "rtDumpSave",
    "dump": {
        "saveName": `${PS4Root}/${name}`,
        "targetFolder": `${PS4Root}/extract/`,
        "selectOnly": [name + '.bin'],
    }
}]
fs.mount("", "/save", {"MountType": "tcpfs"})

for (const cmdObj of cmds) {
    let fh = fs.open("/save/10.0.0.5/1234", 0, 0777)
    let streamWriter = StreamWriter(fh)
    let streamReader = StreamReader(fh)
    let cmd = JSON.stringify(cmdObj);
    streamWriter.writeLine(cmd)
    streamWriter.close()
    let data = streamReader.readLine().trim()
    console.log(data)
    fs.close(fh)
}
let jsonConvert = run("gr2_decode");
let binConvert = run("gr2_encode");
let {execute, convertToSteps} = run("gr2_modder");
let onlinePatchesJson = run("online");
let onlinePatchesSteps = convertToSteps(onlinePatchesJson);

let savePath = `${userPS4Root}/extract/${name}.bin`;
let fi = fs.stat(savePath)
let size = fi.Size()
let fh = fs.open(savePath);
let buff = fs.read(fh, size)
fs.close(fh)
let data = jsonConvert(buff)
execute(data, onlinePatchesSteps)
let result = binConvert(data)

let fmode = fs.constants.O_CREATE 
fmode |= fs.constants.O_WRONLY 
fmode |= fs.constants.O_TRUNC
let fh2 = fs.open("/tmp/" + name + ".bin", fmode, 0777)
fs.write(fh2, result)
fs.close(fh2)


let fh3 = fs.open("/tmp/" + name + ".bin")
discord.uploadFile(name + ".bin", "application/octet-stream", fh3)
fs.close(fh3)

fs.unmount("/save")
