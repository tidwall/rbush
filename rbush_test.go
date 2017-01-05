package rbush

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

type byDims []*nodeT

func (arr byDims) Len() int { return len(arr) }
func (arr byDims) Less(a, b int) bool {
	for i := 0; i < DIMS; i++ {
		n := arr[a].min[i] - arr[b].min[i]
		if n != 0 {
			return n < 0
		}
	}
	for i := 0; i < DIMS; i++ {
		n := arr[a].max[i] - arr[b].max[i]
		if n != 0 {
			return n < 0
		}
	}
	return false
}
func (arr byDims) Swap(a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}

var data = arrToBBoxes(`
	[[0,0,0,0],[10,10,10,10],[20,20,20,20],[25,0,25,0],[35,10,35,10],[45,20,45,20],[0,25,0,25],[10,35,10,35],
    [20,45,20,45],[25,25,25,25],[35,35,35,35],[45,45,45,45],[50,0,50,0],[60,10,60,10],[70,20,70,20],[75,0,75,0],
    [85,10,85,10],[95,20,95,20],[50,25,50,25],[60,35,60,35],[70,45,70,45],[75,25,75,25],[85,35,85,35],[95,45,95,45],
    [0,50,0,50],[10,60,10,60],[20,70,20,70],[25,50,25,50],[35,60,35,60],[45,70,45,70],[0,75,0,75],[10,85,10,85],
    [20,95,20,95],[25,75,25,75],[35,85,35,85],[45,95,45,95],[50,50,50,50],[60,60,60,60],[70,70,70,70],[75,50,75,50],
    [85,60,85,60],[95,70,95,70],[50,75,50,75],[60,85,60,85],[70,95,70,95],[75,75,75,75],[85,85,85,85],[95,95,95,95]]
`)

var emptyData = []*nodeT{
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
	{min: [DIMS]float64{math.Inf(-1), math.Inf(-1)}, max: [DIMS]float64{math.Inf(+1), math.Inf(+1)}},
}

func arrToBBoxes(data string) []*nodeT {
	var nodes []*nodeT
	var arr [][]float64
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		panic(err)
	}
	for _, arr := range arr {
		nodes = append(nodes, &nodeT{
			min: [DIMS]float64{arr[0], arr[1]},
			max: [DIMS]float64{arr[2], arr[3]},
		})
	}
	return nodes
}

func sortedEqual(t *testing.T, a, b []*nodeT) {
	copyA := append([]*nodeT(nil), a...)
	copyB := append([]*nodeT(nil), b...)
	sort.Sort(byDims(copyA))
	sort.Sort(byDims(copyB))
	if !reflect.DeepEqual(copyA, copyB) {
		t.Fatal("not equals")
	}
}

func TestConstructorAcceptsAFormatArgumentToCustomizeTheDataFormat(t *testing.T) {
	// data formats are ignored in go version. test is here for completeness
	tpq("constructor accepts a format argument to customize the data format")
	//var tree = rbush(4, ['.minLng', '.minLat', '.maxLng', '.maxLat']);
	//t.same(tree.toBBox({minLng: 1, minLat: 2, maxLng: 3, maxLat: 4}),
	//{minX: 1, minY: 2, maxX: 3, maxY: 4});
	//t.end();
}
func TestConstructorUses9MaxEntriesByDefault(t *testing.T) {
	tpn("constructor uses 9 max entries by default")
	var tree = New(0).load(someData(9))
	same(t, tree.toJSON().height, 1)

	var tree2 = New(0).load(someData(10))
	same(t, tree2.toJSON().height, 2)
}
func TestToBBoxCompareMinXCompareMinYCanBeOverridenToAllowCustomDataStructures(t *testing.T) {
	// go version uses interfaces instead of callbacks.
	tpq("#toBBox, #compareMinX, #compareMinY can be overriden to allow custom data structures")
	//var tree = rbush(4);
	//tree.toBBox = function (item) {
	//    return {
	//        minX: item.minLng,
	//        minY: item.minLat,
	//        maxX: item.maxLng,
	//        maxY: item.maxLat
	//    };
	//};
	//tree.compareMinX = function (a, b) {
	//    return a.minLng - b.minLng;
	//};
	//tree.compareMinY = function (a, b) {
	//    return a.minLat - b.minLat;
	//};

	//var data = [
	//{minLng: -115, minLat:  45, maxLng: -105, maxLat:  55},
	//{minLng:  105, minLat:  45, maxLng:  115, maxLat:  55},
	//{minLng:  105, minLat: -55, maxLng:  115, maxLat: -45},
	//{minLng: -115, minLat: -55, maxLng: -105, maxLat: -45}
	//];

	//tree.load(data);

	//function byLngLat(a, b) {
	//    return a.minLng - b.minLng || a.minLat - b.minLat;
	//}

	//sortedEqual(t, tree.search({minX: -180, minY: -90, maxX: 180, maxY: 90}), [
	//{minLng: -115, minLat:  45, maxLng: -105, maxLat:  55},
	//{minLng:  105, minLat:  45, maxLng:  115, maxLat:  55},
	//{minLng:  105, minLat: -55, maxLng:  115, maxLat: -45},
	//{minLng: -115, minLat: -55, maxLng: -105, maxLat: -45}
	//], byLngLat);

	//sortedEqual(t, tree.search({minX: -180, minY: -90, maxX: 0, maxY: 90}), [
	//{minLng: -115, minLat:  45, maxLng: -105, maxLat:  55},
	//{minLng: -115, minLat: -55, maxLng: -105, maxLat: -45}
	//], byLngLat);

	//sortedEqual(t, tree.search({minX: 0, minY: -90, maxX: 180, maxY: 90}), [
	//{minLng: 105, minLat:  45, maxLng: 115, maxLat:  55},
	//{minLng: 105, minLat: -55, maxLng: 115, maxLat: -45}
	//], byLngLat);

	//sortedEqual(t, tree.search({minX: -180, minY: 0, maxX: 180, maxY: 90}), [
	//{minLng: -115, minLat: 45, maxLng: -105, maxLat: 55},
	//{minLng:  105, minLat: 45, maxLng:  115, maxLat: 55}
	//], byLngLat);

	//sortedEqual(t, tree.search({minX: -180, minY: -90, maxX: 180, maxY: 0}), [
	//{minLng:  105, minLat: -55, maxLng:  115, maxLat: -45},
	//{minLng: -115, minLat: -55, maxLng: -105, maxLat: -45}
	//], byLngLat);

	//t.end();
}
func TestLoadBulkLoadsTheGivenDataGivenMaxNodeEntir(t *testing.T) {
	tpn("#load bulk-loads the given data given max node entries and forms a proper search tree")
	tree := New(4)
	tree.load(data)
	sortedEqual(t, tree.all(), data)
}
func TestLoadUsesStandardInsertionWhenGivenALowNumberOfItems(t *testing.T) {
	tpn("#load uses standard insertion when given a low number of items")

	var tree = New(8)
	tree.load(data)
	tree.load(data[0:3])

	var tree2 = New(8)
	tree2.load(data)
	tree2.insert(data[0])
	tree2.insert(data[1])
	tree2.insert(data[2])

	same(t, tree, tree2)
}
func TestLoadDoesNothingIfLoadingEmptyData(t *testing.T) {
	tpn("#load does nothing if loading empty data")
	var tree = New(0)
	tree.load(nil)

	same(t, tree, New(0))
}
func TestLoadHandlesTheInsertionOfMaxEntriesPlus2EmptyBBoxes(t *testing.T) {

	tpn("#load handles the insertion of maxEntries + 2 empty bboxes")
	var tree = New(4)
	tree.load(emptyData)

	same(t, tree.data.height, 2)
	sortedEqual(t, tree.all(), emptyData)
}
func TestInsertHandlesTheInsertionOfMaxEntriesPlus2EmptyBBoxes(t *testing.T) {
	tpn("#insert handles the insertion of maxEntries + 2 empty bboxes")
	var tree = New(4)
	var i int
	for _, datum := range emptyData {
		tree.insert(datum)
		i++
	}
	same(t, tree.data.height, 2)
	sortedEqual(t, tree.all(), emptyData)
}
func TestLoadProperlySplitsTreeRootWhenMergingTreesOfTheSameHeight(t *testing.T) {
	tpn("#load properly splits tree root when merging trees of the same height")
	var tree = New(4)
	tree.load(data)
	tree.load(data)
	same(t, tree.data.height, 4)
	sortedEqual(t, tree.all(), append(data, data...))
}

func TestLoadProperlyMergesDataOfSmallerOrBiggerTreeHeights(t *testing.T) {
	tpn("#load properly merges data of smaller or bigger tree heights")
	var smaller = someData(10)

	var tree1 = New(4)
	tree1.load(data)
	tree1.load(smaller)

	var tree2 = New(4)
	tree2.load(smaller)
	tree2.load(data)

	same(t, tree1.data.height, tree2.data.height)

	sortedEqual(t, tree1.all(), append(data, smaller...))
	sortedEqual(t, tree2.all(), append(data, smaller...))

}
func TestSearchFindsMatchingPointsInTheTreeGivenABBox(t *testing.T) {
	tpn("#search finds matching points in the tree given a bbox")

	var tree = New(4)
	tree.load(data)
	var result = tree.search(&nodeT{min: [DIMS]float64{40, 20}, max: [DIMS]float64{80, 70}})

	expect := arrToBBoxes(`[
        [70,20,70,20],[75,25,75,25],[45,45,45,45],[50,50,50,50],[60,60,60,60],[70,70,70,70],
        [45,20,45,20],[45,70,45,70],[75,50,75,50],[50,25,50,25],[60,35,60,35],[70,45,70,45]
	]`)

	sortedEqual(t, result, expect)

}
func TestCollidesReturnsTrueWhenSearchFindsMatchingPoints(t *testing.T) {
	tpn("#collides returns true when search finds matching points")

	var tree = New(4)
	tree.load(data)
	var result = tree.collides(&nodeT{min: [DIMS]float64{40, 20}, max: [DIMS]float64{80, 70}})

	same(t, result, true)
}
func TestSearchReturnsAnEmptyArrayIfNothingFound(t *testing.T) {
	tpn("#search returns an empty array if nothing found")
	var tree = New(4)
	tree.load(data)
	result := tree.search(&nodeT{min: [DIMS]float64{200, 200}, max: [DIMS]float64{210, 210}})

	sortedEqual(t, result, nil)
}

func TestCollidesReturnsFalseIfNothingFound(t *testing.T) {
	tpn("#collides returns false if nothing found")
	var result = New(4).load(data).collides(&nodeT{min: [DIMS]float64{200, 200}, max: [DIMS]float64{210, 210}})

	same(t, result, false)
}
func TestAllReturnsAllPointsInTheTree(t *testing.T) {
	tpn("#all returns all points in the tree")

	var tree = New(4)
	tree.load(data)
	var result = tree.all()

	sortedEqual(t, result, data)
	sortedEqual(t, tree.search(&nodeT{min: [DIMS]float64{0, 0}, max: [DIMS]float64{100, 100}}), data)

}

func TestToJSONAndFromJSONExportsAndImportsSearchTreeInJSONFormat(t *testing.T) {
	// this is not really handling JSON anything, just here for completeness
	tpq("#toJSON & #fromJSON exports and imports search tree in JSON format")

	//var tree = New(4).load(data)
	//var tree2 = New(4).fromJSON(tree.data)

	//sortedEqual(t, tree.all(), tree2.all())
}

func TestInsertAddsAnItemToAnExistingTreeCorrectly(t *testing.T) {
	tpn("#insert adds an item to an existing tree correctly")
	var items = arrToBBoxes(`[
        [0, 0, 0, 0],
        [1, 1, 1, 1],
        [2, 2, 2, 2],
        [3, 3, 3, 3],
        [1, 1, 2, 2]
    ]`)

	var tree = New(4)
	tree.load(items[0:3])
	sortedEqual(t, tree.all(), items[0:3])
	tree.insert(items[3])
	same(t, tree.data.height, 1)
	sortedEqual(t, tree.all(), items[0:4])
	tree.insert(items[4])
	same(t, tree.data.height, 2)
	sortedEqual(t, tree.all(), items)

}
func TestInsertDoesNothingIfGivenUndefined(t *testing.T) {
	tpn("#insert does nothing if given undefined")
	same(t,
		New(0).load(data),
		New(0).load(data).insert(nil))
}
func TestInsertFormsAValidTreeIfItemsAreInsertedOneByOne(t *testing.T) {
	tpn("#insert forms a valid tree if items are inserted one by one")
	var tree = New(4)

	for i := 0; i < len(data); i++ {
		tree.insert(data[i])
	}

	var tree2 = New(4)
	tree2.load(data)

	ok(t, tree.toJSON().height-tree2.toJSON().height <= 1)

	sortedEqual(t, tree.all(), tree2.all())
}

func TestRemoveItemsCorrectly(t *testing.T) {
	tpn("#remove removes items correctly")
	var tree = New(4)
	tree.load(data)

	var len_ = len(data)

	tree.remove(data[0])
	tree.remove(data[1])
	tree.remove(data[2])

	tree.remove(data[len_-1])
	tree.remove(data[len_-2])
	tree.remove(data[len_-3])

	sortedEqual(t,
		data[3:len_-3],
		tree.all())
}
func TestRemoveDoesNothingIfNothingFound(t *testing.T) {
	tpn("#remove does nothing if nothing found")
	same(t,
		New(0).load(data),
		New(0).load(data).remove(&nodeT{min: [DIMS]float64{13, 13}, max: [DIMS]float64{13, 13}}))
}
func TestRemoveDoesIfGivenUndefined(t *testing.T) {
	tpn("#remove does nothing if given undefined")
	same(t,
		New(0).load(data),
		New(0).load(data).remove(nil))
}
func TestRemoveBringsTheTreeToAClearStateWhenRemovingEverythingOneByOne(t *testing.T) {
	tpn("#remove brings the tree to a clear state when removing everything one by one")
	var tree = New(4).load(data)

	for i := 0; i < len(data); i++ {
		tree.remove(data[i])
	}

	same(t, tree.toJSON(), New(4).toJSON())
}
func TestRemoveAcceptsAnEqualsFunction(t *testing.T) {
	// go omits the custom equals function in favor of interfaces
	tpq("#remove accepts an equals function")
	//var tree = New(4).load(data);

	//var item = &nodeT{minX: 20,  70, maxX: 20, maxY: 70, foo: 'bar'};

	//tree.insert(item);
	//tree.remove(JSON.parse(JSON.stringify(item)), function (a, b) {
	//    return a.foo === b.foo;
	//});

	//sortedEqual(t, tree.all(), data);
	//t.end();
}

func TestClearShouldClearAllTheDataInTheTree(t *testing.T) {
	tpn("#clear should clear all the data in the tree")
	same(t,
		New(4).load(data).clear().toJSON(),
		New(4).toJSON())
}
func TestShouldHaveChainableAPI(t *testing.T) {
	tpn("should have chainable API")
	New(0).load(data).insert(data[0]).remove(data[0])
}

func someData(n int) []*nodeT {
	var data []*nodeT
	for i := 0; i < n; i++ {
		data = append(data, &nodeT{min: [DIMS]float64{float64(i), float64(i)}, max: [DIMS]float64{float64(i), float64(i)}})
	}
	return data
}
func ok(t *testing.T, v bool) {
	if !v {
		t.Fatal("not ok")
	}
}
func same(t *testing.T, a interface{}, b interface{}) {
	if a, ok := a.(*RBush); ok {
		if b, ok := b.(*RBush); ok {
			if a.jsonString() != b.jsonString() {
				t.Fatalf("not the same")
			}
			return
		}
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("not the same: '%v' and '%v'", a, b)
	}
}

type byQuick []*nodeT

func (arr byQuick) At(i int) *nodeT {
	return arr[i]
}
func (arr byQuick) Compare(a, b *nodeT) int {
	if a.min[0] < b.min[0] {
		return -1
	}
	if a.min[0] > b.min[0] {
		return 1
	}
	return 0
}
func (arr byQuick) Swap(a, b int) {
	arr[a], arr[b] = arr[b], arr[a]
}

func TestQuickselect(t *testing.T) {
	var arr []*nodeT
	for _, v := range []float64{65, 28, 59, 33, 21, 56, 22, 95, 50, 12, 90, 53, 28, 77, 39} {
		arr = append(arr, &nodeT{min: [DIMS]float64{v}})
	}

	quickselect(arr, 8, 0, len(arr)-1, 0)

	var exp []*nodeT
	for _, v := range []float64{39, 28, 28, 33, 21, 12, 22, 50, 53, 56, 59, 65, 90, 77, 95} {
		exp = append(exp, &nodeT{min: [DIMS]float64{v}})
	}
	for i := 0; i < len(arr); i++ {
		if arr[i].min[0] != exp[i].min[0] {
			t.Fatalf("mismatch for index %d\n", i)
		}
	}
}
func BenchmarkVarious(t *testing.B) {
	t.N = 0
	var N = 1000000
	var maxFill = 16

	fmt.Printf("number: %v\n", N)
	fmt.Printf("maxFill: %v\n", maxFill)

	rand.Seed(time.Now().UnixNano())
	randBox := func(size float64) *nodeT {
		var x = rand.Float64() * (100 - size)
		var y = rand.Float64() * (100 - size)
		return &nodeT{
			min: [DIMS]float64{x, y},
			max: [DIMS]float64{x + size*rand.Float64(), y + size*rand.Float64()},
		}
	}

	genData := func(N int, size float64) []*nodeT {
		var data []*nodeT
		for i := 0; i < N; i++ {
			data = append(data, randBox(size))
		}
		return data
	}
	var start time.Time
	consoleTime := func(s string) {
		start = time.Now()
	}
	consoleTimeEnd := func(s string) {
		end := time.Since(start)
		fmt.Printf("%v: %dms\n", s, end/time.Millisecond)
	}
	var data = genData(N, 1)
	var data2 = genData(N, 1)
	var bboxes100 = genData(1000, 100*math.Sqrt(0.1))
	var bboxes10 = genData(1000, 10)
	var bboxes1 = genData(1000, 1)

	var tree = New(maxFill)

	t.ResetTimer()
	consoleTime("insert one by one")
	for i := 0; i < N; i++ {
		tree.insert(data[i])
		t.N++
	}
	consoleTimeEnd("insert one by one")

	consoleTime("1000 searches 10%")
	for i := 0; i < 1000; i++ {
		tree.search(bboxes100[i])
		t.N++
	}
	consoleTimeEnd("1000 searches 10%")

	consoleTime("1000 searches 1%")
	for i := 0; i < 1000; i++ {
		tree.search(bboxes10[i])
		t.N++
	}
	consoleTimeEnd("1000 searches 1%")

	consoleTime("1000 searches 0.01%")
	for i := 0; i < 1000; i++ {
		tree.search(bboxes1[i])
		t.N++
	}
	consoleTimeEnd("1000 searches 0.01%")

	consoleTime("remove 1000 one by one")
	for i := 0; i < 1000; i++ {
		tree.remove(data[i])
		t.N++
	}
	consoleTimeEnd("remove 1000 one by one")

	consoleTime("bulk-insert 1M more")
	tree.load(data2)
	consoleTimeEnd("bulk-insert 1M more")

	consoleTime("1000 searches 1%")
	for i := 0; i < 1000; i++ {
		tree.search(bboxes10[i])
		t.N++
	}
	consoleTimeEnd("1000 searches 1%")

	consoleTime("1000 searches 0.01%")
	for i := 0; i < 1000; i++ {
		tree.search(bboxes1[i])
		t.N++
	}
	consoleTimeEnd("1000 searches 0.01%")

}

var tpon = false
var tpc int
var tpall string
var tpt int
var tlines []string
var tbad = false
var tbadcount = 0
var tbadidx = 0

func tpsum(s string) string {
	hex := fmt.Sprintf("%X", md5.Sum([]byte(s)))
	return hex[len(hex)-4:]
}
func tpn(format string, args ...interface{}) {
	if tpt == 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("\x1b[92m\x1b[1m✓ %s\x1b[0m\n", fmt.Sprintf(format, args...))
	tpt++
}
func tpq(format string, args ...interface{}) {
	tpn(format, args...)
	return
	if tpt == 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("\x1b[35m\x1b[1m✓ %s\x1b[0m\n", fmt.Sprintf(format, args...))
	tpt++
}

func tpm(format string, args ...interface{}) {
	if tpt == 0 {
		fmt.Printf("\n")
	}
	fmt.Printf("\x1b[34m\x1b[1m• %s\x1b[0m\n", fmt.Sprintf(format, args...))
	tpt++
}
func (this *RBush) jsonString() string {
	var b []byte
	b = append(b, `{`+
		`"maxEntries":`+strconv.FormatInt(int64(this._maxEntries), 10)+`,`+
		`"minEntries":`+strconv.FormatInt(int64(this._minEntries), 10)+`,`+
		`"data":`...)
	b = appendNodeJSON(b, this.data, 1)
	b = append(b, '}')
	return string(b)
}

func appendNodeJSON(b []byte, node *nodeT, depth int) []byte {
	if node == nil {
		return append(b, "null"...)
	}
	b = append(b, '{')
	if len(node.children) > 0 {
		b = append(b, `"children":[`...)
		for i, child := range node.children {
			if i > 0 {
				b = append(b, ',')
			}
			b = appendNodeJSON(b, child, depth+1)
		}
		b = append(b, ']', ',')
	}
	b = append(b, `"leaf":`...)
	if node.leaf {
		b = append(b, "true"...)
	} else {
		b = append(b, "false"...)
	}
	b = append(b, `,"height":`...)
	b = append(b, strconv.FormatInt(int64(node.height), 10)...)
	b = append(b, `,"minX":`...)
	b = append(b, strconv.FormatFloat(node.min[0], 'f', -1, 64)...)
	b = append(b, `,"minY":`...)
	b = append(b, strconv.FormatFloat(node.min[1], 'f', -1, 64)...)
	b = append(b, `,"maxX":`...)
	b = append(b, strconv.FormatFloat(node.max[0], 'f', -1, 64)...)
	b = append(b, `,"maxY":`...)
	b = append(b, strconv.FormatFloat(node.max[1], 'f', -1, 64)...)
	b = append(b, '}')
	return b
}
func nodeJSONString(n *nodeT) string {
	return string(appendNodeJSON([]byte(nil), n, 0))
}