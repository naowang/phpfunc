// phpfunc project phpfunc.go
package phpfunc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesJoin(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func ASerialize(aar []string) []byte {
	rowstr := make([]byte, 0)
	for j := 0; j < len(aar); j++ {
		var wordlen uint32
		wordlen = uint32(len(aar[j]))
		a_buf := bytes.NewBuffer([]byte{})
		binary.Write(a_buf, binary.LittleEndian, wordlen)
		rowstr = BytesJoin(rowstr, a_buf.Bytes(), []byte(aar[j]))
	}
	return rowstr
}

func AASerialize(aar [][]string) []byte {
	totalbin := make([]byte, 0)
	for i := 0; i < len(aar); i++ {
		rowstr := make([]byte, 0)
		for j := 0; j < len(aar[i]); j++ {
			var wordlen uint32
			wordlen = uint32(len(aar[i][j]))
			a_buf := bytes.NewBuffer([]byte{})
			binary.Write(a_buf, binary.LittleEndian, wordlen)
			rowstr = BytesJoin(rowstr, a_buf.Bytes(), []byte(aar[i][j]))
		}
		var sentlen uint32
		sentlen = uint32(len(rowstr))
		a_buf := bytes.NewBuffer([]byte{})
		binary.Write(a_buf, binary.LittleEndian, sentlen)
		totalbin = BytesJoin(totalbin, a_buf.Bytes(), rowstr)
	}
	return totalbin
}

func AAASerialize(aar [][][]string) []byte {
	topbin := make([]byte, 0)
	for i := 0; i < len(aar); i++ {
		totalbin := make([]byte, 0)
		for j := 0; j < len(aar[i]); j++ {
			rowstr := make([]byte, 0)
			for k := 0; k < len(aar[i][j]); k++ {
				var wordlen uint32
				wordlen = uint32(len(aar[i][j][k]))
				a_buf := bytes.NewBuffer([]byte{})
				binary.Write(a_buf, binary.LittleEndian, wordlen)
				rowstr = BytesJoin(rowstr, a_buf.Bytes(), []byte(aar[i][j][k]))
			}
			var sentlen uint32
			sentlen = uint32(len(rowstr))
			a_buf := bytes.NewBuffer([]byte{})
			binary.Write(a_buf, binary.LittleEndian, sentlen)
			totalbin = BytesJoin(totalbin, a_buf.Bytes(), rowstr)
		}
		var sentlen uint32
		sentlen = uint32(len(totalbin))
		a_buf := bytes.NewBuffer([]byte{})
		binary.Write(a_buf, binary.LittleEndian, sentlen)
		topbin = BytesJoin(topbin, a_buf.Bytes(), totalbin)
	}
	return topbin
}

func AUnserialize(mdatastr []byte) []string {
	mdata := make([]string, 0)
	var keylen, startpos int32
	startpos = 0
	for int(startpos)+4 <= len(mdatastr) {
		mdbuf := bytes.NewReader(mdatastr[startpos : startpos+4])
		binary.Read(mdbuf, binary.LittleEndian, &keylen)
		mdata = append(mdata, string(mdatastr[startpos+4:startpos+4+keylen]))
		startpos += 4 + keylen
		if startpos >= int32(len(mdatastr)) {
			break
		}
	}
	return mdata
}

func AAUnserialize(mdatastr []byte) [][]string {
	mdata := make([][]string, 0)
	var keylen, startpos int32
	startpos = 0
	for int(startpos)+4 <= len(mdatastr) {
		mdbuf := bytes.NewReader(mdatastr[startpos : startpos+4])
		binary.Read(mdbuf, binary.LittleEndian, &keylen)
		mdata2str := mdatastr[startpos+4 : startpos+4+keylen]

		mdata2 := make([]string, 0)
		var key2len, startpos2 int32
		startpos2 = 0
		for int(startpos2)+4 <= len(mdata2str) {
			mdbuf2 := bytes.NewReader(mdata2str[startpos2 : startpos2+4])
			binary.Read(mdbuf2, binary.LittleEndian, &key2len)
			mdata2 = append(mdata2, string(mdata2str[startpos2+4:startpos2+4+key2len]))
			startpos2 += 4 + key2len
			if startpos2 >= int32(len(mdata2str)) {
				break
			}
		}
		mdata = append(mdata, mdata2)

		startpos += 4 + keylen
		if startpos >= int32(len(mdatastr)) {
			break
		}
	}
	return mdata
}

func AAAUnserialize(mdatastr []byte) [][][]string {
	mdata := make([][][]string, 0)
	var keylen, startpos int32
	startpos = 0
	for int(startpos)+4 <= len(mdatastr) {
		mdbuf := bytes.NewReader(mdatastr[startpos : startpos+4])
		binary.Read(mdbuf, binary.LittleEndian, &keylen)
		mdata2str := mdatastr[startpos+4 : startpos+4+keylen]

		mdata2 := make([][]string, 0)
		var key2len, startpos2 int32
		startpos2 = 0
		for int(startpos2)+4 <= len(mdata2str) {
			mdbuf2 := bytes.NewReader(mdata2str[startpos2 : startpos2+4])
			binary.Read(mdbuf2, binary.LittleEndian, &key2len)
			mdata3str := mdata2str[startpos2+4 : startpos2+4+key2len]

			mdata3 := make([]string, 0)
			var key3len, startpos3 int32
			startpos3 = 0
			for int(startpos3)+4 <= len(mdata3str) {
				mdbuf3 := bytes.NewReader(mdata3str[startpos3 : startpos3+4])
				binary.Read(mdbuf3, binary.LittleEndian, &key3len)
				mdata3 = append(mdata3, string(mdata3str[startpos3+4:startpos3+4+key3len]))
				startpos3 += 4 + key3len
				if startpos3 >= int32(len(mdata3str)) {
					break
				}
			}
			mdata2 = append(mdata2, mdata3)

			startpos2 += 4 + key2len
			if startpos2 >= int32(len(mdata2str)) {
				break
			}
		}
		mdata = append(mdata, mdata2)

		startpos += 4 + keylen
		if startpos >= int32(len(mdatastr)) {
			break
		}
	}
	return mdata
}

func Getrandmax() int {
	return MaxInt
}

func Rand(minval, maxval int) int {
	i := rand.Intn(maxval - minval)
	return i + minval
}

func RandInt() int {
	return rand.Int()
}

func Array_reverse(ar1 []string) []string {
	ar2 := make([]string, 0)
	for vi := len(ar1) - 1; vi >= 0; vi-- {
		ar2 = append(ar2, ar1[vi])
	}
	return ar2
}

func Preg_match(restr, srcstr string) ([]string, bool) {
	restr2 := restr[1 : len(restr)-1]
	reg := regexp.MustCompile(restr2)
	ar1 := reg.FindSubmatchIndex([]byte(srcstr))
	ar2 := make([]string, 0)
	for i := 0; i < len(ar1); i += 2 {
		foundstr := srcstr[ar1[i]:ar1[i+1]]
		ar2 = append(ar2, foundstr)
	}
	bm := false
	if len(ar2) > 0 {
		bm = true
	}
	return ar2, bm
}

func Str_replace(findstr, withstr, srcstr string) string {
	return strings.Replace(srcstr, findstr, withstr, -1)
}

func Substr(str string, start, size int) string {
	if size > 0 {
		return str[start : start+size]
	} else {
		if len(str)-size >= start {
			return str[start : len(str)+size+1]
		}
	}
	return str
}

func In_array(findval string, ar1 []string) bool {
	for _, val := range ar1 {
		if val == findval {
			return true
		}
	}
	return false
}

func Asort(ar1 []string) []string {
	sort.Strings(ar1)
	return ar1
}

func Array_keys(mdata map[string]string) []string {
	var ldata = make([]string, 0)
	for key, _ := range mdata {
		ldata = append(ldata, key)
	}
	return ldata
}

func Array_search(searchval string, ar1 []string) int {
	for index, val := range ar1 {
		if val == searchval {
			return index
		}
	}
	return -1
}

func ArrayCompare(ar1, ar2 []string) bool {
	if len(ar1) != len(ar2) {
		return false
	}
	for i := 0; i < len(ar1); i++ {
		if ar1[i] != ar2[i] {
			return false
		}
	}
	return true
}

func AArrayCompare(ar1, ar2 [][]string) bool {
	if len(ar1) != len(ar2) {
		return false
	}
	for i := 0; i < len(ar1); i++ {
		if len(ar1[i]) != len(ar2[i]) {
			return false
		}
		for j := 0; j < len(ar1[i]); j++ {
			if ar1[i][j] != ar2[i][j] {
				return false
			}
		}
	}
	return true
}

func AAArrayCompare(ar1, ar2 [][][]string) bool {
	if len(ar1) != len(ar2) {
		return false
	}
	for i := 0; i < len(ar1); i++ {
		if len(ar1[i]) != len(ar2[i]) {
			return false
		}
		for j := 0; j < len(ar1[i]); j++ {
			if len(ar1[i][j]) != len(ar2[i][j]) {
				return false
			}
			for k := 0; k < len(ar1[i][j]); k++ {
				if ar1[i][j][k] != ar2[i][j][k] {
					return false
				}
			}
		}
	}
	return true
}

func Array_key_exists(checkitem string, ar1 []string) bool {
	for _, val := range ar1 {
		if val == checkitem {
			return true
		}
	}
	return false
}

func Str_split(str string, sepsize int) []string {
	strls := make([]string, 0)
	startpos := 0
	for startpos+sepsize <= len(str) {
		strls = append(strls, str[startpos:startpos+sepsize])
		startpos += sepsize
	}
	return strls
}

func Str_V2_Split(str []byte, sepsize int) [][]byte {
	if len(str)%sepsize != 0 {
		panic("str source error!")
	}
	strls := make([][]byte, len(str)/sepsize)
	for i := 0; i < len(str)/sepsize; i++ {
		strls[i] = str[i*sepsize : i*sepsize+sepsize]
	}
	return strls
}

func Array_push(ar1 []string, items ...string) []string {
	for _, item := range items {
		ar1 = append(ar1, item)
	}
	return ar1
}

func EmptyArray() []string {
	itemls := make([]string, 0)
	return itemls
}

func AArray_push(ar1 [][]string, items ...[]string) [][]string {
	for _, item := range items {
		ar1 = append(ar1, item)
	}
	return ar1
}

func EmptyAArray() [][]string {
	itemls := make([][]string, 0)
	return itemls
}

func Array(items ...string) []string {
	itemls := make([]string, 0)
	for _, item := range items {
		itemls = append(itemls, item)
	}
	return itemls
}

func Count(ls []string) int {
	return len(ls)
}

func AACount(ls [][]string) int {
	return len(ls)
}

func AAACount(ls [][][]string) int {
	return len(ls)
}

func Pack(tstr string, num uint64) string {
	if tstr == "V" {
		/*
			var num2 uint32
			num2 = uint32(num)
			a_buf := bytes.NewBuffer([]byte{})
			binary.Write(a_buf, binary.BigEndian, num2)
			return a_buf.String()
		*/
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(num))
		return string(b)
	} else if tstr == "P" {
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, num)
		return string(b)
	}
	return ""
}

func Unpack(tstr, valstr string) uint64 {
	if tstr == "V" {
		valstrbt := []byte(valstr)
		if len(valstrbt) == 4 {
			/*
				mdbuf := bytes.NewReader(valstrbt)
				var valnum uint32
				binary.Read(mdbuf, binary.BigEndian, &valnum)
				return uint64(valnum)
			*/
			return uint64(binary.BigEndian.Uint32([]byte(valstr)))
		}
	} else if tstr == "P" {
		valstrbt := []byte(valstr)
		if len(valstrbt) == 8 {
			return binary.BigEndian.Uint64([]byte(valstr))
		}
	}
	return uint64(0)
}

func Explode(sep, str string) []string {
	ls := make([]string, 0)
	for {
		pos := strings.Index(str, sep)
		if pos == -1 {
			if len(str) > 0 {
				ls = append(ls, str)
			}
			break
		}
		ls = append(ls, str[0:pos])
		str = str[pos+len(sep):]
	}
	return ls
}

func Join(sep string, ls []string) string {
	str := ""
	for i := 0; i < len(ls); i++ {
		if i == 0 {
			str += ls[i]
		} else {
			str += sep + ls[i]
		}
	}
	return str
}

func Array_splice(ls []string, start, cnt int, val []string) []string {
	ls2 := make([]string, 0)
	ls2 = append(ls2, ls[0:start]...)
	ls2 = append(ls2, val...)
	ls2 = append(ls2, ls[start+cnt:]...)

	return ls2
}

func Intval(val float64) int {
	return int(val)
}

func Basename(file_path string) string {
	return filepath.Base(file_path)
}

func File_get_contents(url string) string {

	if strings.Index(url, "http") == 0 {
		return string(file_get_contents_url(url))
	}

	if strings.Index(url, "https") == 0 {
		return string(file_get_contents_url(url))
	}

	return string(file_get_contents_file(url))
	// (strings.Split("a,b,c,d,e,f,g", ",")) // [a b c d e f g]
}

func file_get_contents_url(url string) []byte {

	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func file_get_contents_file(url string) []byte {
	file, err := os.Open(url)
	if err != nil {
		fmt.Println("get content file", url)
		panic(err)
	}
	defer file.Close()

	body, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return body
}

func File_put_contents(fileName string, write_data string) int {
	file, _ := os.Create(fileName)
	defer file.Close()

	wrote_byte, _ := file.Write([]byte(write_data))
	file.Sync()

	return wrote_byte
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func InArray(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
