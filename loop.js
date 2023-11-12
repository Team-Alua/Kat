function walker(root) {
	let dirs = [root];
	while (dirs.length) {
		let root = dirs.shift()
		let data = fs.readdir(root)
		for (const {Name, Dir} of data) {
			let p = root + Name
			if (Dir) {
				p += "/"
				dirs.push(p)
			}
			console.log(p, Dir)
		}
	}
}

fs.mount("/local/abc.zip", "/mnt1", {"MountType": "zipfs", "ReadOnly": true})
// discord.uploadFile("data0001", "application/octet-stream", fh)
// fs.close(fh)
walker("/")

