package templates

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/boltdb/bolt"
)

//line error.ego:1
func Error(w io.Writer, err error) error {
//line error.ego:2
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line error.ego:3
	if _, err := fmt.Fprintf(w, "<!DOCTYPE html>\n"); err != nil {
		return err
	}
//line error.ego:4
	if _, err := fmt.Fprintf(w, "<html lang=\"en\">\n  "); err != nil {
		return err
	}
//line error.ego:5
	if _, err := fmt.Fprintf(w, "<head>\n    "); err != nil {
		return err
	}
//line error.ego:6
	if _, err := fmt.Fprintf(w, "<meta charset=\"utf-8\">\n    "); err != nil {
		return err
	}
//line error.ego:7
	if _, err := fmt.Fprintf(w, "<title>boltd"); err != nil {
		return err
	}
//line error.ego:7
	if _, err := fmt.Fprintf(w, "</title>\n  "); err != nil {
		return err
	}
//line error.ego:8
	if _, err := fmt.Fprintf(w, "</head>\n\n  "); err != nil {
		return err
	}
//line error.ego:10
	if _, err := fmt.Fprintf(w, "<body class=\"error\">\n    "); err != nil {
		return err
	}
//line error.ego:11
	if _, err := fmt.Fprintf(w, "<div class=\"container\">\n      "); err != nil {
		return err
	}
//line error.ego:12
	if _, err := fmt.Fprintf(w, "<div class=\"header\">\n        "); err != nil {
		return err
	}
//line error.ego:13
	if _, err := fmt.Fprintf(w, "<h3 class=\"text-muted\">Error"); err != nil {
		return err
	}
//line error.ego:13
	if _, err := fmt.Fprintf(w, "</h3>\n      "); err != nil {
		return err
	}
//line error.ego:14
	if _, err := fmt.Fprintf(w, "</div>\n\n      An error has occurred: "); err != nil {
		return err
	}
//line error.ego:16
	if _, err := fmt.Fprintf(w, "%v", err); err != nil {
		return err
	}
//line error.ego:17
	if _, err := fmt.Fprintf(w, "\n    "); err != nil {
		return err
	}
//line error.ego:17
	if _, err := fmt.Fprintf(w, "</div> "); err != nil {
		return err
	}
//line error.ego:17
	if _, err := fmt.Fprintf(w, "<!-- /container -->\n  "); err != nil {
		return err
	}
//line error.ego:18
	if _, err := fmt.Fprintf(w, "</body>\n"); err != nil {
		return err
	}
//line error.ego:19
	if _, err := fmt.Fprintf(w, "</html>\n"); err != nil {
		return err
	}
	return nil
}

//line head.ego:1
func head(w io.Writer, tx *bolt.Tx) error {
//line head.ego:2
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line head.ego:4
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
//line head.ego:5
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line head.ego:6
	if _, err := fmt.Fprintf(w, "<head>\n  "); err != nil {
		return err
	}
//line head.ego:7
	if _, err := fmt.Fprintf(w, "<meta charset=\"utf-8\">\n  "); err != nil {
		return err
	}
//line head.ego:8
	if _, err := fmt.Fprintf(w, "<title>"); err != nil {
		return err
	}
//line head.ego:8
	if _, err := fmt.Fprintf(w, "%v", filepath.Base(tx.DB().Path())); err != nil {
		return err
	}
//line head.ego:8
	if _, err := fmt.Fprintf(w, "</title>\n\n  "); err != nil {
		return err
	}
//line head.ego:10
	if _, err := fmt.Fprintf(w, "<style>\n    table {\n      border-collapse:collapse;\n    }\n    \n    table, th, td {\n      border: 1px solid black;\n    }\n\n    th, td { \n      min-width: 100px;\n      padding: 2px 5px;\n    }\n  "); err != nil {
		return err
	}
//line head.ego:23
	if _, err := fmt.Fprintf(w, "</style>\n"); err != nil {
		return err
	}
//line head.ego:24
	if _, err := fmt.Fprintf(w, "</head>\n"); err != nil {
		return err
	}
	return nil
}

//line index.ego:1
func Index(w io.Writer) error {
//line index.ego:2
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line index.ego:3
	if _, err := fmt.Fprintf(w, "<!DOCTYPE html>\n"); err != nil {
		return err
	}
//line index.ego:4
	if _, err := fmt.Fprintf(w, "<html lang=\"en\">\n  "); err != nil {
		return err
	}
//line index.ego:5
	if _, err := fmt.Fprintf(w, "<head>\n    "); err != nil {
		return err
	}
//line index.ego:6
	if _, err := fmt.Fprintf(w, "<meta http-equiv=\"refresh\" content=\"0; url=page\">\n  "); err != nil {
		return err
	}
//line index.ego:7
	if _, err := fmt.Fprintf(w, "</head>\n\n  "); err != nil {
		return err
	}
//line index.ego:9
	if _, err := fmt.Fprintf(w, "<body>redirecting..."); err != nil {
		return err
	}
//line index.ego:9
	if _, err := fmt.Fprintf(w, "</body>\n"); err != nil {
		return err
	}
//line index.ego:10
	if _, err := fmt.Fprintf(w, "</html>\n"); err != nil {
		return err
	}
	return nil
}

//line nav.ego:1
func nav(w io.Writer, tx *bolt.Tx) error {
//line nav.ego:2
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line nav.ego:4
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
//line nav.ego:5
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line nav.ego:6
	if _, err := fmt.Fprintf(w, "<h1>"); err != nil {
		return err
	}
//line nav.ego:6
	if _, err := fmt.Fprintf(w, "%v", filepath.Base(tx.DB().Path())); err != nil {
		return err
	}
//line nav.ego:6
	if _, err := fmt.Fprintf(w, "</h1>\n"); err != nil {
		return err
	}
	return nil
}

//line page.ego:1
func Page(w io.Writer, r *http.Request, tx *bolt.Tx, indexes []int, directID int, showUsage bool) error {
//line page.ego:2
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line page.ego:4
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
//line page.ego:5
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
//line page.ego:6
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}
//line page.ego:7
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line page.ego:9
	p, ids, err := find(tx, directID, indexes)
	if err != nil {
		return err
	}

	// Generate page stats.
	pageSize := tx.DB().Info().PageSize
	stats := p.stats(pageSize)

	// Generate histogram of all nested page usage.
	var histogram map[int]int
	if showUsage {
		histogram = usage(tx, p.id)
	}

//line page.ego:25
	if _, err := fmt.Fprintf(w, "\n\n"); err != nil {
		return err
	}
//line page.ego:26
	if _, err := fmt.Fprintf(w, "<!DOCTYPE html>\n"); err != nil {
		return err
	}
//line page.ego:27
	if _, err := fmt.Fprintf(w, "<html lang=\"en\">\n  "); err != nil {
		return err
	}
//line page.ego:28
	head(w, tx)
//line page.ego:29
	if _, err := fmt.Fprintf(w, "\n\n  "); err != nil {
		return err
	}
//line page.ego:30
	if _, err := fmt.Fprintf(w, "<body>\n    "); err != nil {
		return err
	}
//line page.ego:31
	nav(w, tx)
//line page.ego:32
	if _, err := fmt.Fprintf(w, "\n\n    "); err != nil {
		return err
	}
//line page.ego:33
	if _, err := fmt.Fprintf(w, "<h2>\n      "); err != nil {
		return err
	}
//line page.ego:34
	for i, id := range ids {
//line page.ego:35
		if _, err := fmt.Fprintf(w, "\n        "); err != nil {
			return err
		}
//line page.ego:35
		if i > 0 {
//line page.ego:35
			if _, err := fmt.Fprintf(w, "&raquo;"); err != nil {
				return err
			}
//line page.ego:35
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "\n        "); err != nil {
			return err
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "<a href=\""); err != nil {
			return err
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "%v", pagelink(indexes[:i+1])); err != nil {
			return err
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "\">#"); err != nil {
			return err
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "%v", id); err != nil {
			return err
		}
//line page.ego:36
		if _, err := fmt.Fprintf(w, "</a>\n      "); err != nil {
			return err
		}
//line page.ego:37
	}
//line page.ego:38
	if _, err := fmt.Fprintf(w, "\n    "); err != nil {
		return err
	}
//line page.ego:38
	if _, err := fmt.Fprintf(w, "</h2>\n\n    "); err != nil {
		return err
	}
//line page.ego:40
	if _, err := fmt.Fprintf(w, "<h3>Page Information"); err != nil {
		return err
	}
//line page.ego:40
	if _, err := fmt.Fprintf(w, "</h3>\n    "); err != nil {
		return err
	}
//line page.ego:41
	if _, err := fmt.Fprintf(w, "<p>\n      "); err != nil {
		return err
	}
//line page.ego:42
	if _, err := fmt.Fprintf(w, "<strong>ID:"); err != nil {
		return err
	}
//line page.ego:42
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:42
	if _, err := fmt.Fprintf(w, "%v", comma(int(p.id))); err != nil {
		return err
	}
//line page.ego:42
	if _, err := fmt.Fprintf(w, "<br/>\n      "); err != nil {
		return err
	}
//line page.ego:43
	if _, err := fmt.Fprintf(w, "<strong>Type:"); err != nil {
		return err
	}
//line page.ego:43
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:43
	if _, err := fmt.Fprintf(w, "%v", fmt.Sprintf("%s (%x)", p.typ(), p.flags)); err != nil {
		return err
	}
//line page.ego:43
	if _, err := fmt.Fprintf(w, "<br/>\n      "); err != nil {
		return err
	}
//line page.ego:44
	if _, err := fmt.Fprintf(w, "<strong>Overflow:"); err != nil {
		return err
	}
//line page.ego:44
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:44
	if _, err := fmt.Fprintf(w, "%v", p.overflow); err != nil {
		return err
	}
//line page.ego:44
	if _, err := fmt.Fprintf(w, "<br/>"); err != nil {
		return err
	}
//line page.ego:44
	if _, err := fmt.Fprintf(w, "<br/>\n\n      "); err != nil {
		return err
	}
//line page.ego:46
	if _, err := fmt.Fprintf(w, "<strong>Alloc:"); err != nil {
		return err
	}
//line page.ego:46
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:46
	if _, err := fmt.Fprintf(w, "%v", comma(stats.alloc)); err != nil {
		return err
	}
//line page.ego:46
	if _, err := fmt.Fprintf(w, "<br/>\n      "); err != nil {
		return err
	}
//line page.ego:47
	if _, err := fmt.Fprintf(w, "<strong>In Use:"); err != nil {
		return err
	}
//line page.ego:47
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:47
	if _, err := fmt.Fprintf(w, "%v", comma(stats.inuse)); err != nil {
		return err
	}
//line page.ego:47
	if _, err := fmt.Fprintf(w, "<br/>\n      "); err != nil {
		return err
	}
//line page.ego:48
	if _, err := fmt.Fprintf(w, "<strong>Utilization:"); err != nil {
		return err
	}
//line page.ego:48
	if _, err := fmt.Fprintf(w, "</strong> "); err != nil {
		return err
	}
//line page.ego:48
	if _, err := fmt.Fprintf(w, "%v", fmt.Sprintf("%.2f%%", stats.utilization*100)); err != nil {
		return err
	}
//line page.ego:48
	if _, err := fmt.Fprintf(w, "<br/>\n    "); err != nil {
		return err
	}
//line page.ego:49
	if _, err := fmt.Fprintf(w, "</p>\n\n    "); err != nil {
		return err
	}
//line page.ego:51
	if (p.flags & branchPageFlag) != 0 {
//line page.ego:52
		if _, err := fmt.Fprintf(w, "\n      "); err != nil {
			return err
		}
//line page.ego:52
		if _, err := fmt.Fprintf(w, "<h3>Branch Elements ("); err != nil {
			return err
		}
//line page.ego:52
		if _, err := fmt.Fprintf(w, "%v", p.count); err != nil {
			return err
		}
//line page.ego:52
		if _, err := fmt.Fprintf(w, ")"); err != nil {
			return err
		}
//line page.ego:52
		if _, err := fmt.Fprintf(w, "</h3>\n      "); err != nil {
			return err
		}
//line page.ego:53
		if _, err := fmt.Fprintf(w, "<table>\n        "); err != nil {
			return err
		}
//line page.ego:54
		if _, err := fmt.Fprintf(w, "<thead>\n          "); err != nil {
			return err
		}
//line page.ego:55
		if _, err := fmt.Fprintf(w, "<tr>\n            "); err != nil {
			return err
		}
//line page.ego:56
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Key"); err != nil {
			return err
		}
//line page.ego:56
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:57
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Page ID"); err != nil {
			return err
		}
//line page.ego:57
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:58
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Size (k)"); err != nil {
			return err
		}
//line page.ego:58
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:59
		if _, err := fmt.Fprintf(w, "<th align=\"center\">%%Util"); err != nil {
			return err
		}
//line page.ego:59
		if _, err := fmt.Fprintf(w, "</th>\n          "); err != nil {
			return err
		}
//line page.ego:60
		if _, err := fmt.Fprintf(w, "</tr>\n        "); err != nil {
			return err
		}
//line page.ego:61
		if _, err := fmt.Fprintf(w, "</thead>\n        "); err != nil {
			return err
		}
//line page.ego:62
		if _, err := fmt.Fprintf(w, "<tbody>\n          "); err != nil {
			return err
		}
//line page.ego:63
		for i := uint16(0); i < p.count; i++ {
//line page.ego:64
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:64
			e := p.branchPageElement(i)
//line page.ego:65
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:65
			subpage := pageAt(tx, e.pgid)
//line page.ego:66
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:66
			substats := subpage.stats(pageSize)
//line page.ego:67
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:67
			if _, err := fmt.Fprintf(w, "<tr>\n              "); err != nil {
				return err
			}
//line page.ego:68
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:68
			if _, err := fmt.Fprintf(w, "%v", trunc(tostr(e.key()), 40)); err != nil {
				return err
			}
//line page.ego:68
			if _, err := fmt.Fprintf(w, "</td>\n              "); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "<a href=\""); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "%v", subpagelink(indexes, int(i))); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "\">"); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "%v", e.pgid); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "</a>"); err != nil {
				return err
			}
//line page.ego:69
			if _, err := fmt.Fprintf(w, "</td>\n              "); err != nil {
				return err
			}
//line page.ego:70
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:70
			if _, err := fmt.Fprintf(w, "%v", len(e.key())); err != nil {
				return err
			}
//line page.ego:70
			if _, err := fmt.Fprintf(w, "</td>\n              "); err != nil {
				return err
			}
//line page.ego:71
			if _, err := fmt.Fprintf(w, "<td align=\"right\">"); err != nil {
				return err
			}
//line page.ego:71
			if _, err := fmt.Fprintf(w, "%v", fmt.Sprintf("%.2f%%", substats.utilization*100)); err != nil {
				return err
			}
//line page.ego:71
			if _, err := fmt.Fprintf(w, "</td>\n            "); err != nil {
				return err
			}
//line page.ego:72
			if _, err := fmt.Fprintf(w, "</tr>\n          "); err != nil {
				return err
			}
//line page.ego:73
		}
//line page.ego:74
		if _, err := fmt.Fprintf(w, "\n        "); err != nil {
			return err
		}
//line page.ego:74
		if _, err := fmt.Fprintf(w, "</tbody>\n      "); err != nil {
			return err
		}
//line page.ego:75
		if _, err := fmt.Fprintf(w, "</table>\n    \n    "); err != nil {
			return err
		}
//line page.ego:77
	} else if (p.flags & leafPageFlag) != 0 {
//line page.ego:78
		if _, err := fmt.Fprintf(w, "\n      "); err != nil {
			return err
		}
//line page.ego:78
		if _, err := fmt.Fprintf(w, "<h3>Leaf Elements ("); err != nil {
			return err
		}
//line page.ego:78
		if _, err := fmt.Fprintf(w, "%v", p.count); err != nil {
			return err
		}
//line page.ego:78
		if _, err := fmt.Fprintf(w, ")"); err != nil {
			return err
		}
//line page.ego:78
		if _, err := fmt.Fprintf(w, "</h3>\n      "); err != nil {
			return err
		}
//line page.ego:79
		if _, err := fmt.Fprintf(w, "<table>\n        "); err != nil {
			return err
		}
//line page.ego:80
		if _, err := fmt.Fprintf(w, "<thead>\n          "); err != nil {
			return err
		}
//line page.ego:81
		if _, err := fmt.Fprintf(w, "<tr>\n            "); err != nil {
			return err
		}
//line page.ego:82
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Key"); err != nil {
			return err
		}
//line page.ego:82
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:83
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Value"); err != nil {
			return err
		}
//line page.ego:83
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:84
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Size (k/v)"); err != nil {
			return err
		}
//line page.ego:84
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:85
		if _, err := fmt.Fprintf(w, "<th align=\"center\">%%Util"); err != nil {
			return err
		}
//line page.ego:85
		if _, err := fmt.Fprintf(w, "</th>\n          "); err != nil {
			return err
		}
//line page.ego:86
		if _, err := fmt.Fprintf(w, "</tr>\n        "); err != nil {
			return err
		}
//line page.ego:87
		if _, err := fmt.Fprintf(w, "</thead>\n        "); err != nil {
			return err
		}
//line page.ego:88
		if _, err := fmt.Fprintf(w, "<tbody>\n          "); err != nil {
			return err
		}
//line page.ego:89
		for i := uint16(0); i < p.count; i++ {
//line page.ego:90
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:90
			e := p.leafPageElement(i)
//line page.ego:91
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:91
			if (e.flags & bucketLeafFlag) != 0 {
//line page.ego:92
				if _, err := fmt.Fprintf(w, "\n              "); err != nil {
					return err
				}
//line page.ego:92
				b := ((*bucket)(unsafe.Pointer(&e.value()[0])))
//line page.ego:93
				if _, err := fmt.Fprintf(w, "\n              "); err != nil {
					return err
				}
//line page.ego:94
				util := "-"
				if b.root != 0 {
					substats := pageAt(tx, b.root).stats(pageSize)
					util = fmt.Sprintf("%.2f%%", substats.utilization*100)
				}

//line page.ego:101
				if _, err := fmt.Fprintf(w, "\n              "); err != nil {
					return err
				}
//line page.ego:101
				if _, err := fmt.Fprintf(w, "<tr>\n                "); err != nil {
					return err
				}
//line page.ego:102
				if _, err := fmt.Fprintf(w, "<td>"); err != nil {
					return err
				}
//line page.ego:102
				if _, err := fmt.Fprintf(w, "<strong>"); err != nil {
					return err
				}
//line page.ego:102
				if _, err := fmt.Fprintf(w, "%v", trunc(tostr(e.key()), 40)); err != nil {
					return err
				}
//line page.ego:102
				if _, err := fmt.Fprintf(w, "</strong>"); err != nil {
					return err
				}
//line page.ego:102
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:103
				if _, err := fmt.Fprintf(w, "<td>\n                  &lt;bucket(root="); err != nil {
					return err
				}
//line page.ego:104
				if b.root != 0 {
//line page.ego:104
					if _, err := fmt.Fprintf(w, "<a href=\""); err != nil {
						return err
					}
//line page.ego:104
					if _, err := fmt.Fprintf(w, "%v", subpagelink(indexes, int(i))); err != nil {
						return err
					}
//line page.ego:104
					if _, err := fmt.Fprintf(w, "\">"); err != nil {
						return err
					}
//line page.ego:104
				}
//line page.ego:104
				if _, err := fmt.Fprintf(w, "%v", b.root); err != nil {
					return err
				}
//line page.ego:104
				if b.root != 0 {
//line page.ego:104
					if _, err := fmt.Fprintf(w, "</a>"); err != nil {
						return err
					}
//line page.ego:104
				}
//line page.ego:104
				if _, err := fmt.Fprintf(w, "; seq="); err != nil {
					return err
				}
//line page.ego:104
				if _, err := fmt.Fprintf(w, "%v", b.sequence); err != nil {
					return err
				}
//line page.ego:104
				if _, err := fmt.Fprintf(w, "&gt;\n                "); err != nil {
					return err
				}
//line page.ego:105
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:106
				if _, err := fmt.Fprintf(w, "<td>"); err != nil {
					return err
				}
//line page.ego:106
				if _, err := fmt.Fprintf(w, "%v", len(e.key())); err != nil {
					return err
				}
//line page.ego:106
				if _, err := fmt.Fprintf(w, " / "); err != nil {
					return err
				}
//line page.ego:106
				if _, err := fmt.Fprintf(w, "%v", len(e.value())); err != nil {
					return err
				}
//line page.ego:106
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:107
				if _, err := fmt.Fprintf(w, "<td align=\"right\">"); err != nil {
					return err
				}
//line page.ego:107
				if _, err := fmt.Fprintf(w, "%v", util); err != nil {
					return err
				}
//line page.ego:107
				if _, err := fmt.Fprintf(w, "</td>\n              "); err != nil {
					return err
				}
//line page.ego:108
				if _, err := fmt.Fprintf(w, "</tr>\n            "); err != nil {
					return err
				}
//line page.ego:109
			} else {
//line page.ego:110
				if _, err := fmt.Fprintf(w, "\n              "); err != nil {
					return err
				}
//line page.ego:110
				if _, err := fmt.Fprintf(w, "<tr>\n                "); err != nil {
					return err
				}
//line page.ego:111
				if _, err := fmt.Fprintf(w, "<td>"); err != nil {
					return err
				}
//line page.ego:111
				if _, err := fmt.Fprintf(w, "%v", trunc(tostr(e.key()), 40)); err != nil {
					return err
				}
//line page.ego:111
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:112
				if _, err := fmt.Fprintf(w, "<td>"); err != nil {
					return err
				}
//line page.ego:112
				if _, err := fmt.Fprintf(w, "%v", trunc(tostr(e.value()), 40)); err != nil {
					return err
				}
//line page.ego:112
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:113
				if _, err := fmt.Fprintf(w, "<td>"); err != nil {
					return err
				}
//line page.ego:113
				if _, err := fmt.Fprintf(w, "%v", len(e.key())); err != nil {
					return err
				}
//line page.ego:113
				if _, err := fmt.Fprintf(w, " / "); err != nil {
					return err
				}
//line page.ego:113
				if _, err := fmt.Fprintf(w, "%v", len(e.value())); err != nil {
					return err
				}
//line page.ego:113
				if _, err := fmt.Fprintf(w, "</td>\n                "); err != nil {
					return err
				}
//line page.ego:114
				if _, err := fmt.Fprintf(w, "<td>&nbsp;"); err != nil {
					return err
				}
//line page.ego:114
				if _, err := fmt.Fprintf(w, "</td>\n              "); err != nil {
					return err
				}
//line page.ego:115
				if _, err := fmt.Fprintf(w, "</tr>\n            "); err != nil {
					return err
				}
//line page.ego:116
			}
//line page.ego:117
			if _, err := fmt.Fprintf(w, "\n          "); err != nil {
				return err
			}
//line page.ego:117
		}
//line page.ego:118
		if _, err := fmt.Fprintf(w, "\n        "); err != nil {
			return err
		}
//line page.ego:118
		if _, err := fmt.Fprintf(w, "</tbody>\n      "); err != nil {
			return err
		}
//line page.ego:119
		if _, err := fmt.Fprintf(w, "</table>\n    "); err != nil {
			return err
		}
//line page.ego:120
	}
//line page.ego:121
	if _, err := fmt.Fprintf(w, "\n\n    "); err != nil {
		return err
	}
//line page.ego:122
	if showUsage {
//line page.ego:123
		if _, err := fmt.Fprintf(w, "\n      "); err != nil {
			return err
		}
//line page.ego:124
		mins, maxs, values := bucketize(histogram)
		vmax, maxlen := 0, 20
		for _, v := range values {
			if v > vmax {
				vmax = v
			}
		}

//line page.ego:133
		if _, err := fmt.Fprintf(w, "\n\n      "); err != nil {
			return err
		}
//line page.ego:134
		if _, err := fmt.Fprintf(w, "<h3>Page Usage Histogram"); err != nil {
			return err
		}
//line page.ego:134
		if _, err := fmt.Fprintf(w, "</h3>\n      "); err != nil {
			return err
		}
//line page.ego:135
		if _, err := fmt.Fprintf(w, "<table>\n        "); err != nil {
			return err
		}
//line page.ego:136
		if _, err := fmt.Fprintf(w, "<thead>\n          "); err != nil {
			return err
		}
//line page.ego:137
		if _, err := fmt.Fprintf(w, "<tr>\n            "); err != nil {
			return err
		}
//line page.ego:138
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Usage (bytes)"); err != nil {
			return err
		}
//line page.ego:138
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:139
		if _, err := fmt.Fprintf(w, "<th align=\"left\">Count"); err != nil {
			return err
		}
//line page.ego:139
		if _, err := fmt.Fprintf(w, "</th>\n            "); err != nil {
			return err
		}
//line page.ego:140
		if _, err := fmt.Fprintf(w, "<th>&nbsp;"); err != nil {
			return err
		}
//line page.ego:140
		if _, err := fmt.Fprintf(w, "</th>\n          "); err != nil {
			return err
		}
//line page.ego:141
		if _, err := fmt.Fprintf(w, "</tr>\n        "); err != nil {
			return err
		}
//line page.ego:142
		if _, err := fmt.Fprintf(w, "</thead>\n        "); err != nil {
			return err
		}
//line page.ego:143
		if _, err := fmt.Fprintf(w, "<tbody>\n          "); err != nil {
			return err
		}
//line page.ego:144
		for i := 0; i < len(values); i++ {
//line page.ego:145
			if _, err := fmt.Fprintf(w, "\n            "); err != nil {
				return err
			}
//line page.ego:145
			if _, err := fmt.Fprintf(w, "<tr>\n              "); err != nil {
				return err
			}
//line page.ego:146
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:146
			if _, err := fmt.Fprintf(w, "%v", mins[i]); err != nil {
				return err
			}
//line page.ego:146
			if _, err := fmt.Fprintf(w, " - "); err != nil {
				return err
			}
//line page.ego:146
			if _, err := fmt.Fprintf(w, "%v", maxs[i]); err != nil {
				return err
			}
//line page.ego:146
			if _, err := fmt.Fprintf(w, "</th>\n              "); err != nil {
				return err
			}
//line page.ego:147
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:147
			if _, err := fmt.Fprintf(w, "%v", values[i]); err != nil {
				return err
			}
//line page.ego:147
			if _, err := fmt.Fprintf(w, "</th>\n              "); err != nil {
				return err
			}
//line page.ego:148
			if _, err := fmt.Fprintf(w, "<td>"); err != nil {
				return err
			}
//line page.ego:148
			if _, err := fmt.Fprintf(w, "%v", strings.Repeat("█", int((float64(values[i])/float64(vmax))*float64(maxlen)))); err != nil {
				return err
			}
//line page.ego:148
			if _, err := fmt.Fprintf(w, "</td>\n            "); err != nil {
				return err
			}
//line page.ego:149
			if _, err := fmt.Fprintf(w, "</tr>\n          "); err != nil {
				return err
			}
//line page.ego:150
		}
//line page.ego:151
		if _, err := fmt.Fprintf(w, "\n        "); err != nil {
			return err
		}
//line page.ego:151
		if _, err := fmt.Fprintf(w, "</tbody>\n      "); err != nil {
			return err
		}
//line page.ego:152
		if _, err := fmt.Fprintf(w, "</table>\n    "); err != nil {
			return err
		}
//line page.ego:153
	} else {
//line page.ego:154
		if _, err := fmt.Fprintf(w, "\n      "); err != nil {
			return err
		}
//line page.ego:155
		u, q := r.URL, r.URL.Query()
		q.Set("usage", "true")
		u.RawQuery = q.Encode()

//line page.ego:160
		if _, err := fmt.Fprintf(w, "\n\n      "); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "<p>"); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "<a href=\""); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "%v", u.String()); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "\">Show Page Usage"); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "</a>"); err != nil {
			return err
		}
//line page.ego:161
		if _, err := fmt.Fprintf(w, "</p>\n    "); err != nil {
			return err
		}
//line page.ego:162
	}
//line page.ego:163
	if _, err := fmt.Fprintf(w, "\n\n    "); err != nil {
		return err
	}
//line page.ego:164
	if _, err := fmt.Fprintf(w, "<br/>"); err != nil {
		return err
	}
//line page.ego:164
	if _, err := fmt.Fprintf(w, "<br/>\n    "); err != nil {
		return err
	}
//line page.ego:165
	if _, err := fmt.Fprintf(w, "<form action=\"page\" method=\"GET\">\n      Go to page: "); err != nil {
		return err
	}
//line page.ego:166
	if _, err := fmt.Fprintf(w, "<input type=\"text\" name=\"id\"/>\n      "); err != nil {
		return err
	}
//line page.ego:167
	if _, err := fmt.Fprintf(w, "<button type=\"submit\">Go"); err != nil {
		return err
	}
//line page.ego:167
	if _, err := fmt.Fprintf(w, "</button>\n    "); err != nil {
		return err
	}
//line page.ego:168
	if _, err := fmt.Fprintf(w, "</form>\n  "); err != nil {
		return err
	}
//line page.ego:169
	if _, err := fmt.Fprintf(w, "</body>\n"); err != nil {
		return err
	}
//line page.ego:170
	if _, err := fmt.Fprintf(w, "</html>\n"); err != nil {
		return err
	}
	return nil
}
