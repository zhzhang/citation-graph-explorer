import argparse
import json

def parse_record(record):
    print(record['paperId'])

def load_records(path):
    f = open(path)
    for line in f:
        record = json.loads(line)
        parse_record(record)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='Load the citation graph into memory.',
    )
    parser.add_argument('path')
    args = parser.parse_args()
    path = args.path
    load_records(path)
