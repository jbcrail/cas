#!/usr/bin/env python
# Depends on Python 2.7

import argparse
import errno
import hashlib
import os.path
import sys

from BaseHTTPServer import *

"""
Uses a dictionary-like interface for content-addressable storage.

This allowed me to initially test the server using the in-memory
dictionary instead of the disk-based files.

The directory structure for the files could emulate Git by using
the first two bytes of the SHA-1 as a directory and the remainder
of the SHA-1 as the file containing the content. This would make
it easier to navigate the content for debugging purposes.
"""
class ContentStore(object):
    def __init__(self, path="."):
        if not os.path.exists(path):
            os.makedirs(path)
        self.path = path

    def __contains__(self, key):
        return os.path.exists("%s/%s" % (self.path, key))

    def __getitem__(self, key):
        with open("%s/%s" % (self.path, key)) as f:
            return f.read()

    def __setitem__(self, key, value):
        with open("%s/%s" % (self.path, key), "w") as f:
            f.write(value)

"""
Subclass of default HTTPServer in order to define a default storage
object or allow an already created storage object to be passed in.
"""
class ContentServer(HTTPServer):
    def __init__(self, *args, **kwargs):
        self.store = kwargs.pop('store', ContentStore())
        HTTPServer.__init__(self, *args, **kwargs)

    def serve_forever(self):
        HTTPServer.serve_forever(self)

class ContentRequestHandler(BaseHTTPRequestHandler):
    MIN = 1
    MAX = 64*1024*1024 # 64 MiB

    def response(self, code, mimetype, data):
        self.send_response(code)
        self.send_header('Content-Type', mimetype)
        self.send_header('Content-Length', len(data))
        self.end_headers()
        self.wfile.write(data)

    def do_GET(self):
        """
            200: Content exists, valid SHA-1
            404: Non-existent data requested
            500: Content is corrupted
        """
        sha1 = self.path.lstrip("/")
        if len(sha1) > 0 and sha1 in self.server.store:
            content = self.server.store[sha1]
            content_sha1 = hashlib.sha1(content).hexdigest()
            if content_sha1 == sha1:
                self.response(200, 'application/octet-stream', self.server.store[sha1])
            else:
                self.response(500, 'text/plain', "Content is corrupted")
        else:
            self.response(404, 'text/plain', "SHA1 %s does not exist" % (sha1))

    def do_POST(self):
        """
            201: Previously non-existent content was written successfully
            400: Client content is less than 1 byte or greater than 64MiB
            409: Content already exists on disk
            411: Content-length is required
            500: Generic server error
            507: Disk full
        """
        try:
            value = self.headers.getheader('Content-length', None)
            if value == None:
                self.response(411, 'text/plain', "Content-length header is required")
            elif int(value) < ContentRequestHandler.MIN or int(value) > ContentRequestHandler.MAX:
                self.response(400, 'text/plain', "Content is less than 1 byte or greater than 64 MiB")
            else:
                content = self.rfile.read(int(value))
                sha1 = hashlib.sha1(content).hexdigest()

                if sha1 in self.server.store:
                    self.response(409, 'text/plain', sha1)
                else:
                    self.server.store[sha1] = content
                    self.response(201, 'text/plain', sha1)
        except IOError, e:
            if e == errno.ENOSPC:
                self.response(507, 'text/plain', "Disk full")
            else:
                self.response(500, 'text/plain', "Server error: "+str(e))

if __name__ == '__main__':
    parser = argparse.ArgumentParser(prog="cas")
    parser.add_argument('-p', '--port', type=int, help="port", default=8000)
    parser.add_argument('-d', '--dir', type=str, help="source directory for storage", default=".")
    args = parser.parse_args(sys.argv[1:])
    try:
        server = ContentServer(("",args.port), ContentRequestHandler, store=ContentStore(args.dir))
        server.serve_forever()
    except KeyboardInterrupt:
        server.socket.close()
