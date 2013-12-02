require 'sinatra'
require 'sinatra/partial'
require 'haml'

module CAS
  class Dashboard < Sinatra::Base

    register Sinatra::Partial

    set :haml, :format => :html5
    set :public_folder, File.dirname(__FILE__) + '/public'

    before do
      @services = [
                    ["Go", 8001],
                    ["Python", 8002],
                    ["Ruby", 8003],
                  ]
    end

    get '/' do
      haml :index
    end
  end
end
