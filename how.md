

term21
    global asset engine
        t =themes
        tm=terminal 
        me=model_engine
        while(time+=fps)
            t.time(t) 
            tm.time(t)
            r=me.render([to.obj,tm.obj])
            me.frame(t)


