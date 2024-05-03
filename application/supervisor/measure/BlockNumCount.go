package measure

import (
	"fmt"
	"slices"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/misc"
)

// BlockNumCount to test cross-transaction rate
type BlockNumCount struct {
	epochID int
	bNum    []float64
}

func NewBlockNumCount() *BlockNumCount {
	return &BlockNumCount{
		epochID: -1,
		bNum:    make([]float64, 0),
	}
}

func (ttnc *BlockNumCount) OutputMetricName() string { // Go的语言特性太拉，只能用这个方法了。但你不能加一个Get字样么？GetName不好么。
	return "Tx_number"
}

func (ttnc *BlockNumCount) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // 空的块就算了。
		return
	}
	epochID := b.Epoch
	// extend
	for ttnc.epochID < epochID { // 当然空的块后面也是要补上的。这么弄是为了除掉trailing的空块。但为什么连一句注释都没有？太垃圾了这项目注释写的。
		ttnc.bNum = append(ttnc.bNum, 0)
		ttnc.epochID++
	}

	ttnc.bNum[epochID] += float64(len(b.ExcutedTxs)) // 现在追上了空的块，我们把这个的ExecutedTxs弄上去。
}

func (ttnc *BlockNumCount) OutputRecord() ([]float64, float64) { // 输出的就是这么个（分立统计值，总计值）。问题是统计量都是float64，这合理么？不太合理吧。度量还可以是向量啊甚至是嵌套结构，为何框死成Float64呢。
	return slices.Clone(ttnc.bNum), misc.Sum(ttnc.bNum)
}

func (ttnc *BlockNumCount) GetDesc() metrics.Desc {
	_ = "产生的区块总数计数，单位为 个."
	b := metrics.NewDescBuilder("产生区块数统计(BlockCount)", "对系统产生的区块数量的统计")
	for i, v := range ttnc.bNum {
		b.AddElem(fmt.Sprintf("第%v批次 产生区块数", i+1), "各个批次对共识产生的区块数目的总计", fmt.Sprintf("%v 个", v))
	}
	b.AddElem("合计产生区块数", "整个运行过程中对共识产生的区块数目的总计", fmt.Sprintf("%v 个", misc.Sum(ttnc.bNum)))
	return b.GetDesc()
}
