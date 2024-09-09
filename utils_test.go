package pushbell

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsBase64Key(t *testing.T) {
	hexToBytes := func(h string) []byte {
		data, _ := hex.DecodeString(h)

		return data
	}

	type args struct {
		keyString string
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "AuthParsingPadded",
			args: args{
				"zqyr_PRwQSvvrsytQ9-0YA==",
			},
			want:    hexToBytes("ceacabfcf470412befaeccad43dfb460"),
			wantErr: assert.NoError,
		},
		{
			name: "AuthParsingRaw",
			args: args{
				"gBOK-LZkbIYkL5zU0brnPg",
			},
			want:    hexToBytes("80138af8b6646c86242f9cd4d1bae73e"),
			wantErr: assert.NoError,
		},
		{
			name: "P256dhParsingPadded",
			args: args{
				"BGqW7OPnFArDJMIct3mAC77Eb7PqkgdknuZ2WDCCS5WmLfJopkScXEvla8qxvP19VpONJmHoN1V-00CcTjEM_Lo=",
			},
			want:    hexToBytes("046a96ece3e7140ac324c21cb779800bbec46fb3ea9207649ee6765830824b95a62df268a6449c5c4be56bcab1bcfd7d56938d2661e837557ed3409c4e310cfcba"),
			wantErr: assert.NoError,
		},
		{
			name: "P256dhParsingRaw",
			args: args{
				"BEdvn_bbHUa78RGuOP9M7qv2y9DSAVvdRkV7hbvXEKt1k_ja4_x1JQeMJ5bLji0qDDoF4Zk7qR1RKO9k3ChOV3c",
			},
			want:    hexToBytes("04476f9ff6db1d46bbf111ae38ff4ceeabf6cbd0d2015bdd46457b85bbd710ab7593f8dae3fc7525078c2796cb8e2d2a0c3a05e1993ba91d5128ef64dc284e5777"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsBase64Key(tt.args.keyString)
			if !tt.wantErr(t, err, fmt.Sprintf("parsBase64Key(%v)", tt.args.keyString)) {
				return
			}

			assert.Equal(t, tt.want, result, "should be equal")
		})
	}
}
