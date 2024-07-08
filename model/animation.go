package model

import (
	"github.com/baldurstod/go-source2-tools/kv3"
)

type Animation struct {
	group *AnimGroup
	Name  string
}

func newAnimation(group *AnimGroup) *Animation {
	return &Animation{
		group: group,
	}
}

func (anim *Animation) setData(data kv3.Kv3Value) {
	//log.Println(data)
	/*
		if data == nil {
			return
		}

		animArray, ok := data.GetKv3ValueArrayAttribute("m_animArray")
		if !ok {
			//Should this case be treated as an error ?
			return
		}
		log.Println(animArray)
		/*
			if (data) {
				this.#animArray = data.m_animArray;
				//console.error('data.m_animArray', data.m_animArray);
				this.decoderArray = data.m_decoderArray;
				this.segmentArray = data.m_segmentArray;
				this.frameData = data.m_frameData;

				if (this.#animArray) {
					for (let i = 0; i < this.#animArray.length; i++) {
						let anim = this.#animArray[i];
						this.#animNames.set(anim.m_name, new Source2AnimationDesc(this.animGroup.source2Model, anim, this));
					}
				}
			}
	*/
}
