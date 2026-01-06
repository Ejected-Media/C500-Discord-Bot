# c500-bot-python/fetch_products.py

import requests
import json

# The URL of your local Go API endpoint
API_URL = "http://localhost:8080/api/products"

def get_active_listings():
    """
    Fetches the current product list from the Go backend API.
    """
    print(f"Attempting to connect to: {API_URL}...")

    try:
        # 1. Make the HTTP GET request to the Go server
        response = requests.get(API_URL)

        # 2. Check if the request was successful (status code 200 OK)
        # If it returned a 4xx or 5xx error, this will raise an exception.
        response.raise_for_status()

        # 3. Decode the JSON response into Python data structures (a list of dicts)
        products = response.json()

        print(f"Success! Found {len(products)} products.\n")
        return products

    except requests.exceptions.RequestException as e:
        print(f"Error connecting to API: {e}")
        return []
    except json.JSONDecodeError:
        print("Error: Response was not valid JSON.")
        return []

if __name__ == "__main__":
    # This block runs when you execute the script directly.
    
    listings = get_active_listings()

    if not listings:
        print("No listings found or an error occurred.")
    else:
        print("--- CURRENT ACTIVE LISTINGS ---")
        for item in listings:
            # We can access data using keys representing the JSON structure
            # Notice how the keys match the `json:"..."` tags in your Go models!
            print(f"ID: {item['id']} | Name: {item['name']}")
            print(f"Price: ${item['price']:.2f}")
            print(f"Description: {item['description']}")
            print("-" * 30)

