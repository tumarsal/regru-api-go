package regru

import "context"

// FolderService handles folder-related API calls.
type FolderService struct {
	client *Client
}

// Folder returns the folder service.
func (c *Client) Folder() *FolderService {
	return &FolderService{client: c}
}

// Nop tests folder API availability.
func (s *FolderService) Nop(ctx context.Context) error {
	return s.client.call(ctx, "folder", "nop", nil, nil)
}

// GetList returns the list of folders.
func (s *FolderService) GetList(ctx context.Context) ([]Folder, error) {
	var answer struct {
		Folders []Folder `json:"folders"`
	}
	if err := s.client.call(ctx, "folder", "get_list", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Folders, nil
}

// Create creates a new folder.
func (s *FolderService) Create(ctx context.Context, folderName string) (*Folder, error) {
	params := map[string]interface{}{
		"folder_name": folderName,
	}
	var answer struct {
		FolderID int `json:"folder_id"`
	}
	if err := s.client.call(ctx, "folder", "create", params, &answer); err != nil {
		return nil, err
	}
	return &Folder{
		FolderID:   answer.FolderID,
		FolderName: folderName,
	}, nil
}

// Rename renames a folder.
func (s *FolderService) Rename(ctx context.Context, folderID int, newName string) error {
	params := map[string]interface{}{
		"folder_id": folderID,
		"new_name":  newName,
	}
	return s.client.call(ctx, "folder", "rename", params, nil)
}

// Delete deletes a folder.
func (s *FolderService) Delete(ctx context.Context, folderID int) error {
	params := map[string]interface{}{
		"folder_id": folderID,
	}
	return s.client.call(ctx, "folder", "delete", params, nil)
}

// AddServices adds services to a folder.
func (s *FolderService) AddServices(ctx context.Context, folderID int, services ...ServiceID) error {
	params := map[string]interface{}{
		"folder_id": folderID,
		"services":  services,
	}
	return s.client.call(ctx, "folder", "add_services", params, nil)
}

// RemoveServices removes services from a folder.
func (s *FolderService) RemoveServices(ctx context.Context, folderID int, services ...ServiceID) error {
	params := map[string]interface{}{
		"folder_id": folderID,
		"services":  services,
	}
	return s.client.call(ctx, "folder", "remove_services", params, nil)
}

// ReplaceServices replaces all services in a folder.
func (s *FolderService) ReplaceServices(ctx context.Context, folderID int, services ...ServiceID) error {
	params := map[string]interface{}{
		"folder_id": folderID,
		"services":  services,
	}
	return s.client.call(ctx, "folder", "replace_services", params, nil)
}

// GetServices returns services in a folder.
func (s *FolderService) GetServices(ctx context.Context, folderID int) ([]ServiceListItem, error) {
	params := map[string]interface{}{
		"folder_id": folderID,
	}
	var answer struct {
		Services []ServiceListItem `json:"services"`
	}
	if err := s.client.call(ctx, "folder", "get_services", params, &answer); err != nil {
		return nil, err
	}
	return answer.Services, nil
}

// GetNotInFolderServices returns services not assigned to any folder.
func (s *FolderService) GetNotInFolderServices(ctx context.Context) ([]ServiceListItem, error) {
	var answer struct {
		Services []ServiceListItem `json:"services"`
	}
	if err := s.client.call(ctx, "folder", "get_not_in_folder_services", nil, &answer); err != nil {
		return nil, err
	}
	return answer.Services, nil
}

// MoveServices moves services from one folder to another.
func (s *FolderService) MoveServices(ctx context.Context, fromFolderID, toFolderID int, services ...ServiceID) error {
	params := map[string]interface{}{
		"from_folder_id": fromFolderID,
		"to_folder_id":   toFolderID,
		"services":       services,
	}
	return s.client.call(ctx, "folder", "move_services", params, nil)
}
