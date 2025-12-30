package factories

import (
	"GoHub-Service/app/models/comment"

	"github.com/bxcodec/faker/v3"
)

func MakeComments(count int) []comment.Comment {

	var objs []comment.Comment

	for i := 0; i < count; i++ {
		commentModel := comment.Comment{
			TopicID:   "1",
			UserID:    "1",
			Content:   faker.Sentence(),
			ParentID:  "0",
			LikeCount: 0,
		}
		objs = append(objs, commentModel)
	}

	return objs
}
