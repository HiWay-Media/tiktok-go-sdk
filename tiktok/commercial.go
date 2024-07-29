package tiktok

/*
Query Ads
Use POST /v2/research/adlib/ad/query to query ads.

HTTP URL
https://open.tiktokapis.com/v2/research/adlib/ad/query/

curl -L -X POST 'https://open.tiktokapis.com/v2/research/adlib/ad/query/?fields=ad.id,ad.first_shown_date,ad.last_shown_date' \
-H 'Authorization: Bearer clt.example12345Example12345Example' \
-H 'Content-Type: application/json' \
--data-raw '{
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
}'
*/
func (o *tiktok) ResearchAdQuery(searchTerm string, ) {
    data := map[string]string{}
	data["fields"] = "ad.id,ad.first_shown_date,ad.last_shown_date"
    //
    request := &ResearchAdQueryRequest{
        SearchTerm: searchTerm,
    }
    resp, err := o.restyPostWithQueryParams(API_RESEARCH_AD_QUERY, request, data)
    if err != nil {
        o.debugPrint(err)
		//return nil, err
	}
	/*if resp.IsError() {
		return nil, fmt.Errorf("client access token management error %s", resp.String())
	}*/
	o.debugPrint(resp)
}