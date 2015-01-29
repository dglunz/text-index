# Text Index and Query Server

# API
####File Upload
`POST /index` 

####Search Endpoint
`POST /query?query=`   
search param: `query`

#####Response
Format: `"filename:line-number:word-number"`  
JSON array of occurrences in the previously indexed text  
`{"file1.txt":"125":"5"}` *Line and word counts are 1-based, not 0-based*  
