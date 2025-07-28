#!/usr/bin/env python3
"""
Script to fetch exchange info and convert to CSV format.
Supports multiple exchanges: binance, binanceperp, bybit, etc.
"""

import argparse
import csv
import json
import os
import sys
import requests
from typing import Dict, List, Optional


class ExchangeInfoFetcher:
    """Base class for fetching exchange information."""
    
    def __init__(self, base_url: str):
        self.base_url = base_url
    
    def get_exchange_info(self) -> Dict:
        """Fetch exchange info from the API."""
        raise NotImplementedError("Subclasses must implement get_exchange_info")
    
    def parse_symbols(self, exchange_info: Dict) -> List[Dict]:
        """Parse symbols from exchange info response."""
        raise NotImplementedError("Subclasses must implement parse_symbols")


class BinanceSpotFetcher(ExchangeInfoFetcher):
    """Binance Spot exchange info fetcher."""
    
    def __init__(self):
        super().__init__("https://api.binance.com/api")
    
    def get_exchange_info(self) -> Dict:
        """Fetch exchange info from Binance Spot API."""
        url = f"{self.base_url}/v3/exchangeInfo"
        response = requests.get(url)
        response.raise_for_status()
        return response.json()
    
    def parse_symbols(self, exchange_info: Dict) -> List[Dict]:
        """Parse symbols from Binance Spot exchange info."""
        symbols = []
        for symbol_info in exchange_info.get('symbols', []):
            if symbol_info.get('status') != 'TRADING':
                continue
                
            # Find price filter and lot size filter
            price_filter = None
            qty_filter = None
            
            for filter_info in symbol_info.get('filters', []):
                if filter_info['filterType'] == 'PRICE_FILTER':
                    price_filter = filter_info
                elif filter_info['filterType'] == 'LOT_SIZE':
                    qty_filter = filter_info
            
            if price_filter and qty_filter:
                symbols.append({
                    'symbol': symbol_info['symbol'],
                    'base': symbol_info['baseAsset'],
                    'quote': symbol_info['quoteAsset'],
                    'instrument': 'spot',
                    'priceTick': price_filter.get('tickSize', ''),
                    'szTick': qty_filter.get('stepSize', '')
                })
        
        return symbols


class BinancePerpFetcher(ExchangeInfoFetcher):
    """Binance Perpetual exchange info fetcher."""
    
    def __init__(self):
        super().__init__("https://fapi.binance.com")
    
    def get_exchange_info(self) -> Dict:
        """Fetch exchange info from Binance Perpetual API."""
        url = f"{self.base_url}/fapi/v1/exchangeInfo"
        response = requests.get(url)
        response.raise_for_status()
        return response.json()
    
    def parse_symbols(self, exchange_info: Dict) -> List[Dict]:
        """Parse symbols from Binance Perpetual exchange info."""
        symbols = []
        for symbol_info in exchange_info.get('symbols', []):
            if symbol_info.get('status') != 'TRADING':
                continue
                
            # Find price filter and lot size filter
            price_filter = None
            qty_filter = None
            
            for filter_info in symbol_info.get('filters', []):
                if filter_info['filterType'] == 'PRICE_FILTER':
                    price_filter = filter_info
                elif filter_info['filterType'] == 'LOT_SIZE':
                    qty_filter = filter_info
            
            if price_filter and qty_filter:
                symbols.append({
                    'symbol': symbol_info['symbol'],
                    'base': symbol_info['baseAsset'],
                    'quote': symbol_info['quoteAsset'],
                    'instrument': 'perp',
                    'priceTick': price_filter.get('tickSize', ''),
                    'szTick': qty_filter.get('stepSize', '')
                })
        
        return symbols


class BybitFetcher(ExchangeInfoFetcher):
    """Bybit exchange info fetcher."""
    
    def __init__(self):
        super().__init__("https://api.bybit.com")
    
    def get_exchange_info(self) -> Dict:
        """Fetch exchange info from Bybit API with pagination."""
        all_symbols = []
        cursor = None
        limit = 1000  # Maximum limit per request
        
        # For spot category, no pagination is needed
        categories = ['spot', 'linear', 'inverse']
        
        for category in categories:
            print(f"Fetching {category} instruments...")
            cursor = None
            
            while True:
                url = f"{self.base_url}/v5/market/instruments-info"
                params = {
                    'category': category,
                    'limit': limit
                }
                if cursor:
                    params['cursor'] = cursor
                
                response = requests.get(url, params=params)
                response.raise_for_status()
                data = response.json()
                
                # Add symbols from current page with category info
                page_symbols = data.get('result', {}).get('list', [])
                for symbol in page_symbols:
                    symbol['category'] = category  # Add category info to each symbol
                all_symbols.extend(page_symbols)
                
                # Get next cursor
                next_cursor = data.get('result', {}).get('nextPageCursor')
                
                # If no next cursor or empty page, we're done
                if not next_cursor or not page_symbols:
                    break
                    
                cursor = next_cursor
                print(f"Fetched {len(page_symbols)} {category} symbols, total so far: {len(all_symbols)}")
        
        # Return in the same format as single page response
        return {
            'result': {
                'list': all_symbols
            }
        }
    
    def parse_symbols(self, exchange_info: Dict) -> List[Dict]:
        """Parse symbols from Bybit exchange info."""
        symbols = []
        for symbol_info in exchange_info.get('result', {}).get('list', []):
            if symbol_info.get('status') != 'Trading':
                continue
                
            # Extract price and lot size filters
            price_filter = symbol_info.get('priceFilter', {})
            lot_size_filter = symbol_info.get('lotSizeFilter', {})
            
            # Determine instrument type based on category
            category = symbol_info.get('category', 'spot')
            if category == 'spot':
                instrument = 'spot'
            elif category == 'linear':
                instrument = 'perp'
            elif category == 'inverse':
                instrument = 'inverse'
            else:
                instrument = category  # fallback
            
            # Get szTick - use qtyStep for perpetual, basePrecision for spot
            if category == 'linear' or category == 'inverse':
                sz_tick = lot_size_filter.get('qtyStep', '')
            else:
                sz_tick = lot_size_filter.get('basePrecision', '')
            
            symbols.append({
                'symbol': symbol_info['symbol'],
                'base': symbol_info['baseCoin'],
                'quote': symbol_info['quoteCoin'],
                'instrument': instrument,
                'priceTick': price_filter.get('tickSize', ''),
                'szTick': sz_tick
            })
        
        return symbols


def get_fetcher(market: str) -> ExchangeInfoFetcher:
    """Get the appropriate fetcher for the given market."""
    fetchers = {
        'binance': BinanceSpotFetcher,
        'binanceperp': BinancePerpFetcher,
        'bybit': BybitFetcher,
    }
    
    if market not in fetchers:
        raise ValueError(f"Unsupported market: {market}. Supported markets: {list(fetchers.keys())}")
    
    return fetchers[market]()


def write_csv(symbols: List[Dict], output_path: str):
    """Write symbols to CSV file."""
    fieldnames = ['symbol', 'base', 'quote', 'instrument', 'priceTick', 'szTick']
    
    # Ensure output directory exists
    output_dir = os.path.dirname(output_path)
    if output_dir:
        os.makedirs(output_dir, exist_ok=True)
    
    with open(output_path, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(symbols)
    
    print(f"Successfully wrote {len(symbols)} symbols to {output_path}")


def main():
    parser = argparse.ArgumentParser(
        description="Fetch exchange info and convert to CSV format",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python generate_symbol_table.py binance
  python generate_symbol_table.py binanceperp --dst custom_output.csv
  python generate_symbol_table.py bybit --dst /path/to/output.csv
        """
    )
    
    parser.add_argument(
        'market',
        choices=['binance', 'binanceperp', 'bybit'],
        help='Exchange market to fetch info from'
    )
    
    parser.add_argument(
        '--dst',
        default=None,
        help='Output CSV file path (default: artifact/<market>.csv)'
    )
    
    args = parser.parse_args()
    
    # Set default output path if not provided
    if args.dst is None:
        args.dst = f"artifact/{args.market}.csv"
    
    try:
        # Get the appropriate fetcher
        fetcher = get_fetcher(args.market)
        
        print(f"Fetching exchange info from {args.market}...")
        exchange_info = fetcher.get_exchange_info()
        
        print("Parsing symbols...")
        symbols = fetcher.parse_symbols(exchange_info)
        
        print(f"Found {len(symbols)} trading symbols")
        
        # Write to CSV
        write_csv(symbols, args.dst)
        
    except requests.exceptions.RequestException as e:
        print(f"Error fetching data: {e}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
