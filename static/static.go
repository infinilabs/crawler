package static

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	log "github.com/cihub/seelog"
	"github.com/infinitbyte/framework/core/util"
	"github.com/infinitbyte/framework/core/vfs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

func (vfs StaticFS) prepare(name string) (*vfs.VFile, error) {
	name = path.Clean(name)
	f, present := data[name]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	vfs.once.Do(func() {
		f.FileName = path.Base(name)

		if f.FileSize == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.Compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			log.Error(err)
			return
		}
		f.Data, err = ioutil.ReadAll(gr)

	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return f, nil
}

func (vfs StaticFS) Open(name string) (http.File, error) {

	name = path.Clean(name)

	if vfs.CheckLocalFirst {

		name = util.TrimLeftStr(name, vfs.TrimLeftPath)

		localFile := path.Join(vfs.StaticFolder, name)

		log.Trace("check local file, ", localFile)

		if util.FileExists(localFile) {

			f2, err := os.Open(localFile)
			if err == nil {
				return f2, err
			}
		}

		log.Debug("local file not found,", localFile)
	}

	f, err := vfs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

type StaticFS struct {
	once            sync.Once
	StaticFolder    string
	TrimLeftPath    string
	CheckLocalFirst bool
}

var data = map[string]*vfs.VFile{

	"/assets/css/loadmore.css": {
		FileName:   "assets/css/loadmore.css",
		FileSize:   1721,
		ModifyTime: 1514271253,
		Compressed: `
H4sIAAAAAAAA/5xUsW7bPBDe9RQH/AiQBD8VOY0cm1kK9ElokZIPpkiBYmIlhacuRZfuDVCkQIeOHYp2
aR7HQfMWBU1JluU4dRsggny6j3f3Hb/v5Pg4IO3f8t3dr/v7h9tvy88fHu5+LH++X38jwfFJEIRSM55r
I+B1AAAwR26nFAZRdHCxClhRWcIkZopCIpQV5iJYdHFhoZSo7C78RBsuDIXI/8yZyVBRGMRF1cQ4loVk
1xQmUiez3WVXx7Fklhl9qThJtNSGwn/nfBSNom41YhjHy5LCaVH5eJObpqkPpFpZUuKNoDA4a7L0lTCp
1HMKU+RcKB8tGOeoMgrnvuWnpg85lmwiBa9p0AVL0F5TiMJh3Ie4N2KxKOvkdn6llejVjF1v23BUWR+M
SqISpMPhVGA2tRQGLQ3NdtpAj6/12vyeiPEnDDfzV8RCqSXyDqX1WRNtrc6b7VjDVFkwI5T1WVfCWEyY
bJabI+eynprMxWSGljCFObOoFQWjLbMCovA8LsENyAygSlGhrUF7Jy+C4GVTYSauU8NyUTYYT2Z0UL90
u1mNkGqTNwUOIy6yI199sXrGewIHoy2oo3wv7IvhBtZNs88Uf939v3S9T7c9Z7r9uHz79fHT98c3X550
pjAhTuvPSESX6NeeYiW4D9b3NW7Nx91FCqOoucA3BBUXFYXxeDz2yiqnet4vU6vIfa77CN2TWF1sGt1Z
e3KjtnVk26tMNmGHw7P/ofkP46Mtg+oIqnIWtfKBVlxPCzd+zq7/JLt2aussjCRCyh2jh5hotUnA6Wmf
gHUkRSkpJJfGyf+Vm3CXzy6C3wEAAP//R/EYirkGAAA=
`,
	},

	"/assets/css/m/search_style.css": {
		FileName:   "assets/css/m/search_style.css",
		FileSize:   3422,
		ModifyTime: 1514270523,
		Compressed: `
H4sIAAAAAAAA/+RW227jNhB991cMEAS7C4SGfMkGZV9aFOhvLChybBGhSIGkbGeL9NsLkpItyZLspC8F
SiVIeJkZ8pyZQy5dfiK8dt6U8NcCQEhXKfZGQWolNZJcGf766wKgMk56aTQFi4p5ecAwepTCFxS+Z1l1
Cv0C5b7wFLar1D8W0iNxFeNIQZujZVUYzs2JOPlT6j2F3FiBluQmGuyM9mEKKay2wcf7YtHZ448fwUWF
Nm62ib7Kssdu8NQfGkpd1X7+jOSI+av0xFum29PGrRZMmCMstw6QOXyCnPHXvTW1FufBYP9ZuwQBhezS
IZYJWTsK2wTkxV3YtkMPWfxW1Qke/ogtLjsHoPDwZ2yRPCZEBDvr9IhNcL28pBDtsMJdALFlcA5jgANa
LzlThCm51xRKKYTCEeptydSAYKkLtNJ3kQ/cMst0stHREQCQ0vycmovtem6UfkrbQA6Z5QURyI1lgbEn
uGM5DxEUyWvv77Sw6Grl3WdMLnvrZ+3c+QpzaGpjNmHyTfimnOwMr93oVhkPhR/9m9qH4jkn7Uy0NV+L
tZhMzxtMxS6pFONYGCWa43GjjKXw8Hts005C4nzOlpTu38X+qJmr81ImhbqILcudUbVP8mKqBu6mcuP/
qVw7tVQyu5e6pybj0pLBpjrF36xPDmm3emGuryCNKGx/6av+bVHwePLtKEft0d7WhNqhJQ4Vcn8tClOT
EOibnAxtZHKcEUpz3BmL83dHQr3VVLL9KDLcaI/aU/jyZXIjsboHhdlMdSqT19YF9iojE8AT3mKZX9Xy
RFa6w7575a6f++dr+9Pn20mlZoreYhCNEZ2bLYZV8+xoUH/+nrozFdARoNb/FV7DbP+vpGcLot3n7Gv2
BM3P8vnbBJ53Upywbxk+4z9I7fNTpH0RbAYpvhl7riUtPDAlBfwNtyi/eokxLct4BRLNSqQQ7KMdkbp3
/99Yd+1P1OlupbBcPbtxX8M174vFb62nV3zbWVai68WKx8ke45/Be3JnbEnTE1Exjxvxlayzx0Tht278
5h15x3pTMS79W5Om7wuIOjMdfJhjvUiXybPbVXIbjv3/Ou5i6bhF1K4wKUNLdiJN3r9krcicSJv77dhO
GeZpkqOQLv8EAAD//zt4y59eDQAA
`,
	},

	"/assets/css/m/style.css": {
		FileName:   "assets/css/m/style.css",
		FileSize:   6705,
		ModifyTime: 1514282340,
		Compressed: `
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
eitYQvxPgi/y0ogHD9T/pla4xowTG7k6csnMN9S5akavYMb6aQCK6EB8tjJlVkrQ7bQgDLBNpsQ2oFSR
+4rSC/VDdbnUNK1SKv5igdlTWKVn4RYqn5TglDB8SpdQ+qRSqgjp0px7LuTrKBWCW4pNEko/yteS11UV
rGa9lbBwX1nYh7f/+PPbX375/sXLNy+/jO0Jiz6+AZs6WPxW9bKKWVYc+dlnC8HUZDZ3xbLj0CdBgPkq
6hEYdmTVoQWnNvUCF3Os39Si/CZCK244lZAjTuwKCkPMQ3nbCefIoTfZq47CLjCxogkBDYkE0+YxS0Be
4kCWhPwWQUUiBKrqdS5K5ZQsRRZhnGVIGEdatUirasJt2pnvQ+o2pmT4soACe8FCygwHT9HC5UpXFoqC
qOqpPbSMGaQh3zlh3TuI/awdTNWCnnQziVgYb9hIiKw9rZowjg3rSyN+vMt6XQ6xi215tKf9qYqPCVsK
DrmG1bosZQMQdQHJhiq+GKxdi9ff0SxtLRNyRv0ZRER8o3yM75AKkxN3gYWiVrZp8FRiaPXBK2Hc8try
z0xkR7XetwlrgFbZx/Mt3lO37oxVsSfFR2ObugvPHy/cCF9bXL6bx10Scl2+EESrSC6qlDil7gSxUkUr
Rl/Hm2fyxBpGfKsUR7GZfWOPX69iExMUYhEtZaJSitlHbKdU2aGtBQyvUomHgGFdbOcT4gWUceRzUxn8
9+vf3X319d1fXt59/eo//3z+7vU37148f/vrP7x78VzYF3cJqZi8XGwZeff693ffffvm29/e/fVPG91K
KY09gZKEaqmilN+8fnX31R937OBW+28AAAD//whlhpgxGgAA
`,
	},

	"/assets/css/search/index.css": {
		FileName:   "assets/css/search/index.css",
		FileSize:   209,
		ModifyTime: 1506164087,
		Compressed: `
H4sIAAAAAAAA/0yNzQ6CMAyA73uKJsabJPMgieNpBi1bdSsESpQY3t0o+JOevu9L26g5HeoOZ3gYAIAb
o0YHR2v31VtE4hDVwfkjsh8CiwO7Yu8RWcKXM0uxHSlPtr+vtvbNNQzdJOhg15avqcxizN9n5LFPfnbQ
Jtq2fOIgBSvl0UFDojSs4TKNyu1cNJ0oif7iYp4BAAD//1s9SIXRAAAA
`,
	},

	"/assets/css/search/search.css": {
		FileName:   "assets/css/search/search.css",
		FileSize:   86,
		ModifyTime: 1506166595,
		Compressed: `
H4sIAAAAAAAA/xTGzQmAMAwG0Hun+AbIoeBNp6kmlmK0tT8giLtL3ukF5iqtEVbN23GP3IXAStiTKDfp
tjiqELISCqHYh+J1MGeoMV0zPDym8sAv7nN/AAAA//+r5oC4VgAAAA==
`,
	},

	"/assets/css/search_style.css": {
		FileName:   "assets/css/search_style.css",
		FileSize:   3460,
		ModifyTime: 1527928282,
		Compressed: `
H4sIAAAAAAAA/+RXbW/jNgz+nl9BoCjuDqgC56VXTPuyYcD+xkGWmFioLBmSnKY3dL99kGQljmO7afdl
wJQWrV5Iig/JR8zSlUfCW+dNDX8tAIR0jWKvFKRWUiMpleHPvy4AGuOkl0ZTsKiYlwcMq7XU5EUKX9Ft
UTTHuMSO3RJ8z2vdfFUU92FaodxXnsJ21W1X0iNxDeNIQZsXy5qwXJojcfKn1HsKpbECLSlNFNgZ7cMW
Ulhtg463xaLnyY8fQUWDNro0YTzNh4JSN62fR4K8YPksPfGW6YxJvGrFhHmB5dYBMocPUDL+vLem1eK0
GOQ/K5cgoFCcJ8QyIVtHYZuAPKsL13booYifVXOEuz/iiMdOBijc/RlHDDETIoJd9GbEJrienpKJvKxw
F0BczQY4zw9oveRMEabkXlOopRAKR0Jva6YGAZa6Qit9H/kQW2aZTjI6KgIAUpufU3txXO+Nhp/SbMgh
s7wiArmxLETsAW44zoMFRcrW+xslLLpWefcZkfPdLrN2zr/KHLramE2YchM+U0p2hrdu9KqMB3qI+k3r
Q/GcknbG2pqvxVpMpuc7kYpT0ijGsTJKdO5xo4ylcPd7HNNKQuJ8TpbU7t/Z/qiYa8taJoY6UzIrnVGt
T/Rimg7urnLj/6lce7VUM7uX+oJNxqmlgE1zjL/FZXBIvuo5cpcM0pHC9pfEEbeTgsejz6sctUf7Pie0
Di1xqJD7a1KY2oQQvsnNMEY2xyNCaYk7Y3H+7UioZ04l248iw432qD2FL18mLxKre1CY3VavMnlrXYhe
Y2QCeEJbLPOrWp7ISnfY95/c9eOlf3k+7d9OKjVT9BYDaYzw3GwxrLoupEP98XtuVCYroEdAWf8VXsNs
/6+kZwbR7kv2tXiA7mf5+G0CzxtDnLDPET7hP0jtUyuSO4LNIMU3Y+1a4sIDU1LA3/BeyK86MaZlHZ9A
olmNFIJ8lCNSX7z/75y71ifa9LZSWK4e3biu4Zm3xeK3rOkZX3eW1egubEV3ivv4Z9BP7oytaWoRFfO4
EV/JurhPIfzWt9/1kTecNw3j0r92afq2gMgz08aHOXZh6bx5UrtKaoPb/y93F0vHLaJ2lUkZ2vvq89T7
NpRzP6/tlGGeJjoK6fJPAAAA//8L3V9jhA0AAA==
`,
	},

	"/assets/css/style.css": {
		FileName:   "assets/css/style.css",
		FileSize:   7061,
		ModifyTime: 1520233310,
		Compressed: `
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
AA==
`,
	},

	"/assets/css/tasks.css": {
		FileName:   "assets/css/tasks.css",
		FileSize:   143,
		ModifyTime: 1496674736,
		Compressed: `
H4sIAAAAAAAA/zTLUQrCMAzG8feeIheoB2gOI2kTRfphRjYpIru7uNWXkO8H/8tNHjC9brL2TyIiqtL6
Pfz11EJvA3zw4c3hUShMOe1pkwqjWXioRW4OyLJa+T9nNzw01zDp5bhZAD51hCwTf5vT/g0AAP//ivV6
Eo8AAAA=
`,
	},

	"/assets/js/footer.js": {
		FileName:   "assets/js/footer.js",
		FileSize:   2572,
		ModifyTime: 1514271647,
		Compressed: `
H4sIAAAAAAAA/5RV227bRhO+N+B3GDM/QtKRSMv+g7pUGUNWjCJF6qaxjRZwDWO1HEpbkbv07lKHxrov
UPR9+kDtexRLUTwoqltTF9rDzPfNef3Dw/09OIShRKIxgtESvkF9LgnjCj5MsistZOoVIjcKZQDfYkST
Yv+WaAygd9ztHXf/X5xcsxQD+POP3/76/dfToHdkDv39PfPz/a8ObodvB9eDW7NlMThWziOMGcfIgoMQ
9DJDEUNKFhkZowuf9vcAAGZEQq5QDsbINYTAyYyNiRbSq049Ld6LOcohUei4/VqPqXuRoSQQ1hAe4xEu
vosdu7iyXcPd7cHLl1AceDOUigm+DZSKXyAEp6bPpIhyqiEMwf4a6VTYrgGpiVQ+Ulo6O5hjJjEWC9uF
V3DagZMtKoaGaYdeqhg2DT7YOPifmUuEV/C6pF0BJgqbwS4TACEcmfv9vUjQPDUogk9xmWcQQpxzqpng
4NSJMhoIZ4AQwJzxSMw9nCHXDd80GXOSGsG1l2eAnpL0IsECX5PxpbkOAD1N5BiroxLDVE2FEYL97vLD
zbUNj4/QPL2++PF68PFiYLsgUeeSl9qE6pwkQxEVlnpTXBbrs8baMNMJkWbT4GxqhnB6Uvlsvv859ouH
HOXSdr1Y0FxVlbN9qzBBqqvrlWmEDcUm6m+g10L/nP7ky5bARmh3N3Fc6PtWPzW/Mk2JoKTIZljL99vS
q3q7etK4L55hXCZx9hzjKvl/Na4I7apvxs7d3ZudE+jn701WnLqSWxOH9BplK/gg14KKNEtQ41WRxVYT
zEiSYwcioslWaVgvYss1/ZiyRuL7m8wXVJ9h/8AiPYEQqsYboy57RJ0vhwlRyjSFY6nRoktzpUV6fz+X
JMtQWq7HNKbOkeuJOFZYgnXhuOGQysdjVPo+lwmEhZUKiaSTLhVc40IDuaWGJbQJ1WyG9p3lekRr6VgN
Vas5tjaaz4QrJSooUywt60KoSsf0+dadZbXi3fbL8su9b1VVUaU0M5lTELbUUc4YxRuZBE2oTi0xN8EM
/illDcEIE5YybR5M3+k89t2f1KHfuBd8rRfsLK4WUIzyIz7kqPT5MoDXRx3w/ZQlTCEVPFK1aEYkSVUA
n8oXLAC75x3ZnaLdArAYz3JtdTapCjaLVYOOiyGhEwxAyxzB9xVq0KLYdcwiYoqMEgRK6ITxcV3P5p/0
ivTXA68ZJ6cMuVsU/8p9ojmrui9w1iER0rG9dVmNxAJujU+hJVGhtu4MVRRdmNfmPVMaOUrHpgmjU7sD
OzpcT5jyMiKR60sR4TZPEahiXuuB1pKNco2OVTS51TE1138WTP0oPOl1Zed6WlwVvran0nS+eyqcL99F
jr2JemFoaSIVXIkEvUSMnencNZRp9uBpSejUHKxf+NKkvwMAAP//Wmg4PAwKAAA=
`,
	},

	"/assets/js/page/console.js": {
		FileName:   "assets/js/page/console.js",
		FileSize:   5350,
		ModifyTime: 1518269207,
		Compressed: `
H4sIAAAAAAAA/7RYX2/rthV/96fg5QJYmh0pfegenChDm5vuZsh6uyR361AUAS0dy0woUiEpO1ma7z6Q
+kfRcuZeoHoIIvLwd8758fyTj4JVxVNNBUdBiF4nE4QQ2hCJUsH5afdWqBwl6CjAfypUjsN+g4l2gwl/
I6fc2TRvrkBBnq9Fnnx7clor7ewgZQk8u4GnCpS+FnlQqNyYhppHwlONepbRDUoZUSqZylp8en5WSjjH
aGZNniF8FpuFszijm/NWvXlqLQZewlNol988Q0oiFfxUqTVkH4kmvh2iTAqVR6pkVAcY4fk3DjxdBaL8
kFQ8gxXlkIX9OfMUKq8PV0ulZeDj/HLya8SA53o9+2ZuNusXB7/TkST4p5urf313d4k9He9wpUrBFRxO
Vg+nK8nHuGufN2AKesu+fH99dfHHG+bbZWLtAPMuPv/4w9XfxsxLBVeCQcREHmDKqTYrK5rjcEdULB+S
v99+/jGy4WKDZJc4A9HYFYjlwz4PhrZOdv/z7NoQRjO0ErIgeoGwCZU2lM1fJvJb+l9ITk7HMuxSSiEP
zS8wwtPzs+XwVpZfl1m9iKebrgIm8ihdU5ZJ4EHYRP55UpcL76qOjHQYrSjPgmlGN4sVlUpPw0hCITYQ
eDw3fBwfn+5Sa7KsNutOWNRepjk2m526zjQlzSU4ju8+f/yMSpD2SngKiBalFBvY434bp6MstOjv+W5k
fpf/LegoBzskWPgBEe15jwxT443PODQlraDaayu9Zyj4YFqLu4j6DF4RpmDMMnvQmLchLAgPPe0y2h02
iX99i3/7zV1JmdpZASK9tYvry+9u/HrhBew49V5+O6YZNiIFPOst3EmmYTNshAa3Zpawm4UqlYKxj3Tj
WuLyVF9feNq3fAkrCWr9Q8VTlHTRGjjuDsqPBGM6pNpp7K1DkOpW75ujQRdeJepkXx1CDPFXXIM0Tumi
NbENgy3lmdj+gv8Ny1uRPoLGv/rxsBZKowQxkRKjJjLvnBSAZijoFkshNformi6maIaGqws0nXoXGMdb
pSZ+/+2PSaFFKliS4LXWpVrs6SocJYjDFnXGB3ir1CKOTVm1Zs8QjrfKb3K2cR2OeADg0BcbhYKLErhz
9yiAjfbJRf1MmOCLJgQER6A0WTJqxqUI77bAg9IENbWkiYt7pYmuFA4jDc96T289CnBUPR4TBlLjFvjC
tK0At+vHGeE5SBxGJMv8PVWlKagdhtB4LI5cgZsTRnmbSUhLmlu1Q+pPx6hPmVDwtdzbw6O0/yFsdoyN
0dlSvYuti0SB7uh06s0cfXtycvL/2FWgv4raApQi+UHk+jM/bHSUEU12sqf7D9mRcojkDVj4P6KSaCnF
VoFEmQCFuNBIVaWtNl3mqshl7c1tsE5Rndii/fDPCuRLkIm0KoBrc1Mkewn8qn3kSCztmPAIL1UZqenc
IcOpVd03XrQSaaWCsOsTVm13xp1oLUGv3SDARH7PYAMMh7Yvme2orNT6vtuJtPhSliAviIKu39FVLWp0
3JdEa5D8Q4L/jD3z3H1XhbseDsaTDpkyeBfZ2R8gO+vjyE2IvQfuibj43lboz1YmI2kq+HHJyIvJOq1l
gDOqyJJBhucWRAJhmhYQ7pxSWpQjpz54x96c663KjGhoLxi4PdJ2aeNNLXBffxbdt19ubW2xo4iV1c86
OTLfKmWlG5fdmcOYkLw2aTuwJ2mUOlvDGErGYi10xN14SPYEzkDeueVkTzi48t6tJfsv2R1h4pgJkqFM
FIRyC9RHSUQeyHMwrCWVZAs0jRVoM2jFZgYHOZ0PZPRLCQuES6E0Hu4Y/AWy36dKS8pzunqp83Uol5J0
DWhRj4a7EHe1ggcluKegaQWLfl408nNkwuDW9ps5enj6+dPNWLX9cvVIdcSFNlbh274P72tL7yVBG6Tv
nt2XCt3hYZkf+mo/gh1Pf/7H9Sety2Y6H/psZe/WUmx3vnV2Pb+sgc3I5p7bO7HZUmzydTKJ43oejgS3
cTUc3OMYNYmmUILc5oudDwQ8R385qftvHL+5Vd4Reu1ytv9pr+0tUQ76koH59/uXqyzA3u99zWtUw92J
0g7n7tonoPla+78rvofvYPu4Hubb5H8BAAD///LvfQvmFAAA
`,
	},

	"/assets/js/page/index.js": {
		FileName:   "assets/js/page/index.js",
		FileSize:   1072,
		ModifyTime: 1496674736,
		Compressed: `
H4sIAAAAAAAA/4xTz2rbThC+71MM+v3AkmNWii8BG9NDC22vrSmFEsxaGsVrr1fO7qyTNBh6DvTYngpt
XyDkbvI2jZvHKFrF/4gC0UGH+b6Z/Wbmm7jZZNCElwYFYQbDC5hiliooNLSTw6P4MD7iDJoxy51OSRYa
rMjxNdIHoRyG8+iSAQDIPJz3ek5nmEuN2UO0/AySMxqSro8s2E5s3mULti2sCpG9EiTC6JJ5WhzD3dWv
+9vbv8uf99e/V99uVl+vPfI/P0EKG7ElQbYR8azQGG4KhZkgEcFWg8x96FNw6tAhT0eYToLjHZG+Zhj8
5xE0AxJ2MtBuGkSc8JzqsvmsmB0EEENwUIs6OzoIWjAXSmadDenhAe6MCo65BwcWMYu6Gy0L9lh3asSZ
QsNncoZKaqwXX5Eeid9b2FP1eC61tCPMor2maohUkFBla2hMYcrWnveAp0dl4tCgmDw/0dN9oh/xE4nV
+HOkdFQtJ6rZzgZ2drQ3cv+PumxtvMp3qx9fVt9v/iyXd1f77hNjcR5uN+CM6kAjLgf/wsrP2GsnjdYG
pYtZKfoEKdgGS1X9ChjbQu8g1qUpWtuBtZ99gy0od/meBDnbgvHpxzfvdh3uSxpx1hdDVU2EG7ROUX2X
C8ZITrEv0wkS9MAivdWEZi5UuD5DaLWTJCnZ28Pssn8BAAD//w3N87swBAAA
`,
	},

	"/assets/js/page/tasks.js": {
		FileName:   "assets/js/page/tasks.js",
		FileSize:   532,
		ModifyTime: 1515907893,
		Compressed: `
H4sIAAAAAAAA/2SRzWrjMBSF93qKQwhIcoJ/MjAb402GySyGQAvdtCELxZYTNbYVZCWt2+Tdi2U7LlSL
D6F77tH9CTyPwMMfI4WVGXYNSpmlBXSFRRj9DqJFsAh9Ai8gJD9XqVW6glWlFHvNLvyTAIA1TXdpj5H2
bCrg9fEsTeOPWqe4pcKmByb5j4RLFye3bx/l2pTCLhsra7ZrOc9kqkpR1Bydgcq7AJIEIR+8aIg2icZO
cxEGRySIwjAcX7ISCQY7zBDhesWvMV6rj9YWG+r+p3PQ/8uWa8d/jk+OD45/HV8cn5d0O1opJFgLe/Dz
QmvD3LXQ+65yjgD3lyPnXVrfyEmYWq4KLSzr++zFJ/3GjnMozn2rV+pdZiwrOccMFBSzrvqN2sbtPKds
mCjr5z5lk2EzE37fEY/JjcfkKwAA//9Y+UgoFAIAAA==
`,
	},

	"/assets/js/snapshot_footprint.js": {
		FileName:   "assets/js/snapshot_footprint.js",
		FileSize:   254,
		ModifyTime: 1509258526,
		Compressed: `
H4sIAAAAAAAA/1yNMUvEQBBG+/0VH6k2BHL2cREUCxsL0Uoslt1JMrCZkb2Jp8j9dwl3itxUA2/emwNL
1kOvUjRmBIyrJGMV3+LbOQD4iBVR0qx1j4CsaV1IrJ/I7gtt6/726zlOj3Eh38SmHU7aqBV+cxkBVwMY
17+ZvpBMNg/grtve4Dw8en8+eeW33ujT2vaP/iNzpREBza5S5krJdjdrLaFBB5KkmV6eHu50eVchMX8h
nopHd3TO/QQAAP//nDYkC/4AAAA=
`,
	},

	"/favicon.ico": {
		FileName:   "favicon.ico",
		FileSize:   9662,
		ModifyTime: 1514805677,
		Compressed: `
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
//+jkx0iviUAAA==
`,
	},

	"/assets": {
		IsFolder: true,
		FileName: "/assets",
	},

	"/assets/css": {
		IsFolder: true,
		FileName: "/assets/css",
	},

	"/assets/css/m": {
		IsFolder: true,
		FileName: "/assets/css/m",
	},

	"/assets/css/search": {
		IsFolder: true,
		FileName: "/assets/css/search",
	},

	"/assets/js": {
		IsFolder: true,
		FileName: "/assets/js",
	},

	"/assets/js/page": {
		IsFolder: true,
		FileName: "/assets/js/page",
	},
}
