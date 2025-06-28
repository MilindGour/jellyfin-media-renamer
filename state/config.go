package state

import "github.com/MilindGour/jellyfin-media-renamer/models"

// LastConfigSourceById stores the result of DirectoryEntry of last call to /api/config/source/id
var LastConfigSourceById []models.DirectoryEntry = nil

// LastSecondPageAPIResponse stores the result of latest copy of SecondScreenResponse
var LastSecondPageAPIResponse *models.SecondScreenResponse = nil
