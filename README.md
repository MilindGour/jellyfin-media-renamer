# Requirements of the project

## First points
- Configurable options
  - List of download directories (Can be more than one). Title also can be configured
    - Name
    - Path
  - List of Jellyfin media directories (Can be more than one). Title also can be configured

## Screen flows
### First Screen
  - Dropdown menu with all the download directories in it.
  - User selects one of the directories and presses *Scan* button.
  - System will fetch all the subdirectories of the selected directory.
  - Subdirectories will be displayed along with their size
  - User selects one / many items from the list.
  - Only the selected items need to be processed in the next screen.
  - User presses *Next* button

### Second Screen
  - User will see all the selected items from first screen and a *Get IDs* button
  - A drop-down menu with all the media database providers (like TMDB, IMDB etc)
  - Each item contains following additional fields:
    - Name, size
    - Type of media as multiple buttons (TV, Movie) -> preselected on some logic (find if TV or Movie)
    - User can change media type if system detected wrong
  - User selects a media db provider and presses *Get IDs* button
  - System will do following:
    - Clean the filenames to make it searchable
    - Fetches using the cleaned names (can use info like year in tmdb -> movieName y:1992)
    - Appends the search results to their respective list items (logic below)
  - If a media item has:
    - 0 result: prompt user to adjust search string and search that item again.
    - 1 result: select the id as the id of the result
    - 2 or more results: user sees list of results (name, year, description) and let user select one result.
  - After all IDs are determined, enable *Next* button.
  - User presses *Next* button

### Third Screen
  - Now we have selected media items, their IDs and media type
