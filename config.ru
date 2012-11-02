require ::File.expand_path('../app',  __FILE__)

log = ::File.new('logs/sinatra.log', 'a')
STDOUT.reopen(log)
STDERR.reopen(log)

run Sinatra::Application
