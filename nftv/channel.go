package nftv

type channel int

const (
	channelNFTV1 channel = iota
	channelNFTV2
	channelNFTV3
)

var channelIDs = []string{
	channelNFTV1: "NFTV-1.us",
	channelNFTV2: "NFTV-2.us",
	channelNFTV3: "NFTV-3.us",
}

var channelIDMap = map[string]channel{
	channelIDs[channelNFTV1]: channelNFTV1,
	channelIDs[channelNFTV2]: channelNFTV2,
	channelIDs[channelNFTV3]: channelNFTV3,
}

var channelNames = []string{
	channelNFTV1: "NFTV 1",
	channelNFTV2: "NFTV 2",
	channelNFTV3: "NFTV 3",
}

var channelURLs = []string{
	channelNFTV1: "https://www.nightflightplus.com/guide/nftv-1",
	channelNFTV2: "https://www.nightflightplus.com/guide/nftv-2",
	channelNFTV3: "https://www.nightflightplus.com/guide/nftv-3",
}

var channelIcons = []string{
	channelNFTV1: "https://image.c.cdn.zype.com/5527e30469702d5e08000000/64f8f4ab3c85bb0001ad6507/custom_thumbnail/240.jpg",
	channelNFTV2: "https://image.c.cdn.zype.com/5527e30469702d5e08000000/63ee80d96d8ecd0001c332f3/custom_thumbnail/240.jpg",
	channelNFTV3: "https://image.c.cdn.zype.com/5527e30469702d5e08000000/67c218e0978fc00001d39ea7/custom_thumbnail/240.jpg",
}

func (c channel) id() string {
	return channelIDs[c]
}

func (c channel) name() string {
	return channelNames[c]
}

func (c channel) url() string {
	return channelURLs[c]
}

func (c channel) icon() string {
	return channelIcons[c]
}

func getChannelByID(id string) channel {
	return channelIDMap[id]
}
