// Package state contains variables which are used
// to store temporary values for dependent APIs.
package state

import "github.com/MilindGour/jellyfin-media-renamer/models"

// LastConfigSourceByID stores the result of DirectoryEntry of last call to /api/config/source/id
var LastConfigSourceByID *models.ConfigSourceByIDResponse = nil

// LastSecondPageAPIResponse stores the result of latest copy of SecondScreenResponse
var LastSecondPageAPIResponse *models.SecondScreenResponse = nil

var LastConfirmedMediaIds *models.ScrapSearchRequest = nil

var LastScrapConfirmRequest *models.ScrapSearchConfirmRequest = nil
