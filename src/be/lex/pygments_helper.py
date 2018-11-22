# -*- coding:utf-8 -*-

import json
import sys

import pygments
import pygments.lexers
from pygments.lexers import get_lexer_for_filename


def parse(filename):
    lexer = get_lexer_for_filename(filename, stripnl=False, stripall=False)
    f = open(filename, 'r', encoding='utf-8')
    content = f.read()
    f.close()
    tokens = list(lexer.get_tokens(content))
    result = []
    line = []
    beg_pos = 0
    for token in tokens:
        token_type = token[0]
        token_content = token[1]
        token_parts = token_content.split(u"\n")
        if len(token_parts) == 1:
            # 没有出现换行
            target = {
                "content": token_content,
                "type": token_type,
                "beg_pos": beg_pos,
                "end_pos": beg_pos + len(token_content)
            }
            beg_pos += len(token_content)
            line.append(target)
        else:
            for idx, token_part in enumerate(token_parts):
                if idx != len(token_parts)-1:
                    target = {
                        "content": token_part,
                        "type": token_type,
                        "beg_pos": beg_pos,
                        "end_pos": beg_pos + len(token_part)
                    }
                    if target['beg_pos'] != target['end_pos']:
                        line.append(target)
                    result.append(line)
                    line = []
                    beg_pos = 0
                else:
                    target = {
                        "content": token_part,
                        "type": token_type,
                        "beg_pos": beg_pos,
                        "end_pos": beg_pos + len(token_part)
                    }
                    if target['beg_pos'] != target['end_pos']:
                        line.append(target)
                    beg_pos += len(token_part)
    if len(line) != 0:
        result.append(line)
    print(json.dumps(result))
        
            
if __name__ == "__main__":
    try:
        parse(sys.argv[1])
    except Exception:
        raise
        print("ERROR")