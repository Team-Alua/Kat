while(true) {
	let data = discord.getMessage();
	if (data.Content == "upload") {
		discord.sendMessage("I bet you want to upload something");
	} else if (data.Content == "end") {
		break
	} else {
		discord.sendMessage("I don't understand " + data.Content)
	}
}
