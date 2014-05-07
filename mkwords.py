#!/usr/bin/env python
import random
import sys

default = """
V:a/e/i/o/u
C:p/t/c/s/k
#CVC
CV(C/CC/CV/)
"""

DROPOFF = 0.7


def main():
    rules = Rules(*sys.argv[1:])
    print "Generating {} word{}".format(
        rules.num, 's' if rules.num > 1 else '')
    for i in range(0, rules.num):
        print rules.word(),
    print


class Rules(object):
    def __init__(self, source=None, num=1):
        if source in (None, '-', ''):
            source = default
        else:
            source = open(source).read()
        self.source = source
        self.classes = {}
        self.patterns = []
        self.num = int(num)
        self.parse()

    def word(self):
        letters = []
        # print self.patterns
        pat = random.choice(self.patterns)
        # print pat
        for tok in tokens(pat):
            letters.append(self.choose(tok))
        return ''.join(letters)

    def parse(self):
        lines = self.source.split("\n")
        for line in lines:
            line = line.strip()
            if line.startswith("#") or not line:
                continue
            if ':' in line:
                # print 'class line', line
                try:
                    cls, rule = line.split(':')
                except ValueError:
                    pass
                else:
                    choices = weight(rule.split('/'))
                    self.classes.setdefault(cls.strip(), []).append(
                        choices)
                    # print 'added', cls, rule
            else:
                # print 'pattern line', line
                self.patterns.append(line)

    def choose(self, cls):
        # print cls, self.classes[cls]
        choices = random.choice(self.classes[cls])
        return weighted_choice(choices)


def weighted_choice(choices):
    total = sum(w for c, w in choices)
    r = random.uniform(0, total)
    upto = 0
    for c, w in choices:
        if upto + w > r:
            return c
        upto += w
    assert False, "Shouldn't get here"


def tokens(rule):
    parts = iter(rule)
    while True:
        for item in next_token(parts):
            yield item


def next_token(gen):
    t = gen.next()
    lt = []
    if t == '(':
        weight = 100.0
        choices = []
        lt = []
        while True:
            t = next_token(gen)
            if t == ')':
                choices.append((lt, weight))
                return weighted_choice(choices)
            elif t == '/':
                choices.append((lt, weight))
                weight = weight * DROPOFF
                lt = []
            else:
                lt.append(t)
    else:
        return t

def weight(choices):
    out = []
    weight = 100.0
    for item in choices:
        out.append((item, weight))
        weight = weight * DROPOFF
    return out


if __name__ == '__main__':
    main()
