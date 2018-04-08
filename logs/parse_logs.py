#! /usr/bin/env python

import argparse, re

class MongoDBLogger():
    list = []
    def process(self, args):
        filename = args.file
        pattern = '^(\d+-\d+-\d+T\d+:\d+:\d+.\d+)(\S+) (.?) (\S+)\s+\[\S+\] (.*) (\d+)ms$'
        pattern = '^(\d+-\d+-\d+T\d+:\d+:\d+.\d+)(\S+) (.?) (\S+)\s+\[\S+\] command (\S+) command: (\S+) (.*) (\d+)ms$'
        queryPat = '^(.*) filter: ({.*? }), (.*) planSummary: (\S+) (.*)'

        with open(filename) as f:
            lines = f.readlines()
            for message in lines:
                match = re.match(pattern, message)
                if not match:
                    continue
                ns = match.group(5)
                cmd = match.group(6)
                text = match.group(7)
                ms = int(match.group(8))
                #print(cmd)
                # print('%s %s %s %d' % (scan, ns, cmd, ms))
                if cmd in ['find', 'getMore']:
                    qm = re.match(queryPat, text)
                    if not qm:
                        #print(text)
                        continue
                    else:
                        filter = qm.group(2)
                        scan = qm.group(4)
                        # str = '%s %s %s %s %d ms' % (ns, cmd, scan, filter, ms)
                        #if str not in self.list:
                        #    self.list.append(str)
                        doc = {'ns': ns, 'cmd': cmd, 'scan': scan, 'filter': filter, 'ms': ms}
                        self.list.append(doc)
                        # print('* %s %s %s' % (cmd, filter, scan))
                elif cmd in ['replSetHeartbeat', 'replSetUpdatePosition', 'isMaster']:
                    continue
                elif cmd in ['insert']:
                    continue
                elif cmd in ['update']:
                    continue
                else:
                    continue

        sort_on = 'ms'
        decorated = [(dict_[sort_on], dict_) for dict_ in self.list]
        decorated.sort(reverse=True)
        mlist = decorated[:args.topn]
        for item in mlist:
            print(item)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("-f", "--file", default='mongod.log', help="input file")
    parser.add_argument("-n", "--topn", default=20, help="input file")
    args = parser.parse_args()

    if not args.file:
        print('--file is required')
        exit()

    m = MongoDBLogger()
    m.process(args)
