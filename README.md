# nf-scraper
A simple API made to scrape and return a brasilian digital receipt in JSON
WARNING: Currently only works with Minas Gerais receipts 

## How to run
 - `docker build --tag nf-scraper .`
 - `docker run -p 8000:8000 nf-scraper`