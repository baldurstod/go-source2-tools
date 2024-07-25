package particles

import "github.com/baldurstod/go-source2-tools/kv3"

type PreviewState struct {
	ModSpecificData   int
	SequenceName      string
	HitboxSetName     string
	MaterialGroupName string
	PreviewModel      string
}

func (ps *PreviewState) initFromDatas(datas *kv3.Kv3Element) error {
	ps.ModSpecificData, _ = datas.GetIntAttribute("m_nModSpecificData")
	ps.SequenceName, _ = datas.GetStringAttribute("m_sequenceName")
	ps.HitboxSetName, _ = datas.GetStringAttribute("m_hitboxSetName")
	ps.MaterialGroupName, _ = datas.GetStringAttribute("m_materialGroupName")
	ps.PreviewModel, _ = datas.GetStringAttribute("m_previewModel")
	return nil
}
