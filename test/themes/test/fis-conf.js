fis.config.merge({
    namespace: 'test',
    modules : {
        parser : {
            less: 'less'
        },
        postprocessor: {
            js: 'jswrapper'
        }
    },
    roadmap : {
        ext : {
            less: 'css'
        },
        path : [
            {
                reg : /^\/page\/(.+\.html)$/i,
                isMod: true,
                url : '${namespace}/page/$1',
                release : '/template/${namespace}/page/$1',
                extras: {
                    isPage: true
                }
            },
            {
                reg: /^\/(static|test)\/(.*)/i,
                release: '/$1/${namespace}/$2'
            },
            {
                reg: '${namespace}-map.json',
                release: '/config/${namespace}-map.json'
            },
            {
                reg: /^.+$/,
                release: '/static/${namespace}/$&'
            }
        ]
    },
    settings : {
        postprocessor : {
            jswrapper: {
                type: 'amd'
            }
        }
    }
});
