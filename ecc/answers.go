package ecc

func h160ToP2pkhAddress(h160 []byte, testnet bool) {
	var prefix []byte
	testnet = false
	if bool(testnet) {
		prefix = []byte("o")
	} else {
		prefix = []byte("\u0000")
	}
	prefix = append(prefix, h160...)
	encodeBase58Checksum(string(prefix))
}

func h160ToP2shAddress(h160 []byte, testnet bool) {
	var prefix []byte
	testnet = false

	if bool(testnet) {
		//b'\xc4'
		prefix = []byte("o")
	} else {
		//b'\x05'
		prefix = []byte("\u0000")
	}
	prefix = append(prefix, h160...)
	encodeBase58Checksum(string(prefix))
}
