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
fs.umount("/mnt1")
// discord.uploadFile("data0001", "application/octet-stream", fh)
// fs.close(fh)
walker("/")

