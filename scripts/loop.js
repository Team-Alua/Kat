function walker(root) {
	let dirs = [root];
	while (dirs.length) {
		let root = dirs.shift()
		let data = fs.readdir(root)
		for (const {Name, Dir, Size} of data) {
			let p = root + Name
			if (Dir) {
				p += "/"
				dirs.push(p)
			}
			console.log(p, Dir, Size)
		}
	}
}

fs.mount("/local/abc.zip", "/mnt1", {"MountType": "zipfs", "ReadOnly": true})

fs.mount("/tmp/abc.zip", "/mnt2", {"MountType": "zipfs", "ReadOnly": false})
fs.copyDir("/mnt1", "/mnt2")
fs.unmount("/mnt2")

fs.unmount("/mnt1")


let fh = fs.open("/tmp/abc.zip")
discord.uploadFile("upload.zip", "application/zip", fh)
fs.close(fh)
// walker("/")

