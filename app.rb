require 'rubygems'
require 'net/http'
require 'bundler/setup'
Bundler.require(:default, ENV['RACK_ENV'].to_sym)

configure :production, :development do
  set :server, :puma
end

get '/:job.png' do
  send_file JenkinsStatus.new(params[:job]).image
end

get '/:job' do
  "[![Build Status](https://ci-status.moozement.net/#{params[:job]}.png)](http://ci.moozement.net/job/#{params[:job]}/)"
end

class JenkinsStatus

  attr_reader :job, :color

  def initialize(job)
    @job = job
    @color = 'gray'
    fetch
  end

  def image
    "public/#{status}.png"
  end

  private

  def fetch
    uri = URI("http://ci.moozement.net/job/#{job}/api/json")
    req = Net::HTTP::Get.new(uri.request_uri)
    res = Net::HTTP.start(uri.host, uri.port) { |http| http.request(req) }
    @color = JSON.parse(res.body)['color'] if res.is_a?(Net::HTTPSuccess)

  end

  def status
    case color
      when 'blue' then 'passing'
      when 'red' then 'failing'
      else 'unknown'
    end
  end
end
