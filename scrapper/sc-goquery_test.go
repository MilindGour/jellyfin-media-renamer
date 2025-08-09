package scrapper

import (
	"testing"

	"github.com/MilindGour/jellyfin-media-renamer/network"
)

func TestGoQuery_splitAttribute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		selector string
		want     string
		want2    string
	}{
		{
			name:     "Selector with attribute present",
			selector: ".parent .item a.link[href]",
			want:     ".parent .item a.link",
			want2:    "href",
		},
		{
			name:     "Selector without attribute",
			selector: ".parent .item .title",
			want:     ".parent .item .title",
			want2:    "",
		},
		{
			name:     "Selector with attribute at wrong place",
			selector: ".parent[data-id] .title",
			want:     ".parent[data-id] .title",
			want2:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g GoQuery
			got, got2 := g.splitAttribute(tt.selector)
			if got != tt.want {
				t.Errorf("element splitAttribute() = <%v>, want <%v>", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("attribute splitAttribute() = <%v>, want <%v>", got2, tt.want2)
			}
		})
	}
}

func TestGoQuery_Scrap(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		url      string
		itemSel  string
		fieldMap map[string]string
		want     ScrapResultList
		wantErr  bool
	}{
		{
			name:     "Basic scrapping",
			url:      "mock-scrap-html",
			itemSel:  ".scrap1",
			fieldMap: map[string]string{"test_val": ".target-val", "test_attr": "a.target-attr[href]"},
			want:     ScrapResultList{{"test_val": "Test Value", "test_attr": "Test Attr"}},
			wantErr:  false,
		},
		{
			name:     "Advanced scrapping",
			url:      "mock-scrap-html",
			itemSel:  ".item-list .item",
			fieldMap: map[string]string{"title": "span.title", "img": "img[src]"},
			want:     ScrapResultList{{"title": "Title 1", "img": "src_1"}, {"title": "Title 2", "img": "src_2"}, {"title": "Title 3", "img": "src_3"}},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g GoQuery = GoQuery{
				htmlProvider: network.NewMockHtml(),
			}
			got, gotErr := g.Scrap(tt.url, tt.itemSel, tt.fieldMap)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Scrap() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Scrap() succeeded unexpectedly")
			}

			// make sure they are same size
			if len(tt.want) != len(got) {
				t.Errorf("Scrap() = %d items, wanted %d items", len(got), len(tt.want))
			}
			for i, ttWant := range tt.want {
				// compare input and output maps
				for ttKey, ttValue := range ttWant {
					gotVal, gotKey := got[i][ttKey]
					if !gotKey {
						t.Errorf("Scrap() = key '%s' not found.", ttKey)
					}
					if gotVal != ttValue {
						t.Errorf("Scrap() = %v, want %v", gotVal, ttValue)
					}
				}
			}
		})
	}
}
