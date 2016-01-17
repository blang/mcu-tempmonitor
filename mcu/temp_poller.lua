
function start_temp_poller() 
    sk = nil
    connected=0
    --sk=net.createConnection(net.UDP, 0)
    --sk:connect(10001,"192.168.0.131")
    tmr.alarm(1,2000,1,function()
        print("Connected",connected)
        if (connected == 1) then
            temp=getTemp(cfg.sensors.temp)
            print("Temp:")
            print(temp)
            sk:send(temp)
        end
        if (connected == 0) then 
            sk=net.createConnection(net.UDP, 0)
            print("Connect")
            print(cfg.target.host)
            sk:connect(cfg.target.port,cfg.target.host)
    
            connected=1
        end
        end)
    tmr.start(1)
end