package impl

import (
	"Hermes/internal/models"
	"context"
)

func (s *Service) GetComments(ctx context.Context, params models.QueryParams) ([]models.Comment, error) {

	roots, err := s.storage.GetRootComments(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []models.Comment

	for _, root := range roots {

		flat, err := s.storage.GetCommentTree(ctx, root.ID)
		if err != nil {
			return nil, err
		}

		tree := buildTree(flat)
		if len(tree) > 0 {
			result = append(result, *tree[0])
		}

	}

	return result, nil

}

func buildTree(comments []models.Comment) []*models.Comment {

	hm := make(map[int64]*models.Comment)
	var roots []*models.Comment

	for i := range comments {
		hm[comments[i].ID] = &comments[i]
	}

	for _, c := range comments {

		node := hm[c.ID]
		if node.ParentID == nil {
			roots = append(roots, node)

		} else {

			parent, ok := hm[*node.ParentID]
			if ok {
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}

		}
	}

	return roots

}
