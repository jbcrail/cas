require 'sinatra'
require 'digest/sha1'

helpers do
    def exists?(id)
        File.exists? id
    end

    def retrieve(id)
        File.open(id) { |f| f.read }
    end

    def store(id, content)
        File.open(id, 'w') { |f| f.write(content) }
    end

    def hash(content)
        Digest::SHA1.hexdigest content
    end
end

before do
    content_type 'text/plain'
end

configure do
    mime_type :binary, 'application/octet-stream'
end

get '/:id' do |id|
    if not exists? id
        halt 404, "SHA-1 #{id} does not exist"
    end

    content = retrieve(id)
    sha1 = hash(content)
    if sha1 == id
        status 200
        content_type :binary
        body content
    else
        status 500
        body "Content is corrupted"
    end
end

post '/' do
    content_length = request.content_length.to_i
    if content_length < 1 or content_length > 64*1024*1024
        halt 400, "Content is less than 1 byte or greater than 64MiB"
    end

    request.body.rewind
    content = request.body.read

    id = hash(content)
    if exists? id
        status 409
    else
        status 201
        begin
            store(id, content)
        rescue Errno::ENOSPC
            halt 507, "Disk full"
        rescue Exception => msg
            halt 500, "Server error: #{msg}"
        end
    end
    body id
end
