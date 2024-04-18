package pushbell

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsBase64Key(t *testing.T) {
	a := assert.New(t)

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
			name: "auth parsing padded",
			args: args{
				"zqyr_PRwQSvvrsytQ9-0YA==",
			},
			want:    []byte{206, 172, 171, 252, 244, 112, 65, 43, 239, 174, 204, 173, 67, 223, 180, 96},
			wantErr: assert.NoError,
		},
		{
			name: "auth parsing raw",
			args: args{
				"gBOK-LZkbIYkL5zU0brnPg",
			},
			want:    []byte{128, 19, 138, 248, 182, 100, 108, 134, 36, 47, 156, 212, 209, 186, 231, 62},
			wantErr: assert.NoError,
		},
		{
			name: "p256dh parsing padded",
			args: args{
				"BGqW7OPnFArDJMIct3mAC77Eb7PqkgdknuZ2WDCCS5WmLfJopkScXEvla8qxvP19VpONJmHoN1V-00CcTjEM_Lo=",
			},
			want:    []byte{4, 106, 150, 236, 227, 231, 20, 10, 195, 36, 194, 28, 183, 121, 128, 11, 190, 196, 111, 179, 234, 146, 7, 100, 158, 230, 118, 88, 48, 130, 75, 149, 166, 45, 242, 104, 166, 68, 156, 92, 75, 229, 107, 202, 177, 188, 253, 125, 86, 147, 141, 38, 97, 232, 55, 85, 126, 211, 64, 156, 78, 49, 12, 252, 186},
			wantErr: assert.NoError,
		},
		{
			name: "p256dh parsing raw",
			args: args{
				"BEdvn_bbHUa78RGuOP9M7qv2y9DSAVvdRkV7hbvXEKt1k_ja4_x1JQeMJ5bLji0qDDoF4Zk7qR1RKO9k3ChOV3c",
			},
			want:    []byte{4, 71, 111, 159, 246, 219, 29, 70, 187, 241, 17, 174, 56, 255, 76, 238, 171, 246, 203, 208, 210, 1, 91, 221, 70, 69, 123, 133, 187, 215, 16, 171, 117, 147, 248, 218, 227, 252, 117, 37, 7, 140, 39, 150, 203, 142, 45, 42, 12, 58, 5, 225, 153, 59, 169, 29, 81, 40, 239, 100, 220, 40, 78, 87, 119},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsBase64Key(tt.args.keyString)
			if !tt.wantErr(t, err, fmt.Sprintf("parsBase64Key(%v)", tt.args.keyString)) {
				return
			}

			a.Equal(tt.want, got, "should be equal")
		})
	}
}

func Test_parseApplicationKeys(t *testing.T) {
	k1 := "Oow2uOCuGBpkGupoKT7R3vrM25Zd8L6mg0fXPgxFPqY"
	p1 := "BIIBOPkUO37N3K0ibzYGJw0ZkGS-dDTjVZSCpdrhTZBntuK_aAyLN_nlbgGz6fkjwcu6cmFkPhOaqm5YiGeo8Y0"

	k2 := "Oow2uOCuGBpkGupoKT7R3vrM25Zd8L6mg0fXPgxFPqY="
	p2 := "BIIBOPkUO37N3K0ibzYGJw0ZkGS-dDTjVZSCpdrhTZBntuK_aAyLN_nlbgGz6fkjwcu6cmFkPhOaqm5YiGeo8Y0="

	type args struct {
		publicKey  *string
		privateKey *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "keys non-padded",
			args: args{
				publicKey:  &p1,
				privateKey: &k1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "keys padded",
			args: args{
				publicKey:  &p2,
				privateKey: &k2,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseApplicationKeys(tt.args.publicKey, tt.args.privateKey)
			if !tt.wantErr(t, err, fmt.Sprintf("parseApplicationKeys(%v, %v)", tt.args.publicKey, tt.args.privateKey)) {
				return
			}
		})
	}
}
