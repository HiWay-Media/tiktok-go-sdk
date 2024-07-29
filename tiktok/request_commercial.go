package tiktok

/*{
   "filters":{
       "advertiser_business_ids": [3847236290405, 319282903829],
       "ad_published_date_range": {
            "min": "20210102",
            "max": "20210109"
       },
       "country_code": "FR",
       "unique_users_seen_size_range": {
           "min": "10K",
           "max": "1M"
       },
   },
   "search_term": "mobile games"
}*/

type ResearchAdQueryRequest struct {
    Filters ResearchAdQueryFilter `json:"filters"`
    SearchTerm string `json:"search_term"`
}

type ResearchAdQueryFilter struct {
    AdvertiserBusinessIDs   []int64                 `json:"advertiser_business_ids"`
    AdPublishedDateRange    AdPublishedDateRange    `json:"ad_published_date_range"`
    CountryCode             string                  `json:"country_code"`
}

type AdPublishedDateRange struct {
    Min string `json:"min"`
    Max string `json:"max"`
}