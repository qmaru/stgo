import React, { useState, useEffect } from 'react'
import { makeStyles } from '@material-ui/styles'

import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import Grid from '@material-ui/core/Grid'
import Snackbar from '@material-ui/core/Snackbar'
import Button from '@material-ui/core/Button'
import Skeleton from '@material-ui/lab/Skeleton'
import LazyLoad from 'react-lazyload'
import Box from '@material-ui/core/Box'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'

const useStyles = makeStyles(theme => ({
    mainBar: {
        position: 'fixed',
        backgroundColor: "#dd4b9a"
    },
    mainGrid: {
        position: 'relative',
        top: 100
    },
    stBody: {
        margin: "0 auto"
    },
    stCard: {
        margin: 10,
        textAlign: "center",
    },
    stCardContent: {
        paddingBottom: 0
    },
    stTitle: {
        padding: 10,
        textAlign: "left"
    },
    stImg: {
        width: "100%"
    },
    stDown: {
        padding: 10
    },
    stDownCls: {
        color: "#fff",
        backgroundColor: "#f491c7",
        '&:hover': {
            backgroundColor: "#dd4b9a",
        }
    },
    stLoad: {
        padding: 10,
        height: 200,
    },
    SkeletonCls: {
        margin: '10px auto'
    }
}))


export default function Index() {
    const classes = useStyles()
    const [stInfo, setSTInfo] = useState([])
    const [minID, setMinID] = useState(9999)
    const [message, setMessage] = useState("")

    const [open, setOpen] = useState(false)
    const [loading, setLoading] = useState(false)
    const [stLoadState, setSTLoadState] = useState(false)
    const [stLoadHidden, setSTLoadHidden] = useState({ display: "none" })
    const [stLoadText, setSTLoadText] = useState("more")

    const getStList = () => {
        setSTLoadState(true)
        let url = ""
        if (minID === 9999) {
            url = `http://localhost:61800/stlist`
        } else {
            url = `http://localhost:61800/stlist?next_id=${minID}`
        }
        fetch(url,
            {
                method: "GET",
                dataType: "json",
            })
            .then(res => res.json())
            .then(response => {
                if (response.status === 1) {
                    let sids = []
                    let data = stInfo
                    let add_data = response.data
                    for (let i in add_data) {
                        let d = add_data[i]
                        sids.push(parseInt(d["ID"]))
                    }
                    let minid = Math.min.apply(null, sids)
                    let newdata = data.concat(add_data)
                    setSTInfo(newdata)
                    setMinID(minid)
                    setSTLoadHidden({ display: "block" })
                    setSTLoadState(false)
                    setLoading(false)
                } else {
                    setSTLoadHidden({ display: "block" })
                    setSTLoadState(true)
                    setSTLoadText("no data")
                    setLoading(false)
                }
            }).catch(
                () => {
                    setOpen(true)
                    setLoading(false)
                    setSTLoadState(false)
                    setMessage("Server Error")
                }
            )
    }

    useEffect(() => {
        getStList()
    }, [])

    const barClose = (event, reason) => {
        if (reason === 'clickaway') {
            return
        }
        setOpen(false)
    }

    const Loading = () => {
        return (
            <Box width="65%" >
                <Skeleton width="24%" height={60} className={classes.SkeletonCls} />
                <Skeleton width="100%" height={30} />
                <Skeleton width="100%" height={30} />
                <Skeleton width="80%" height={30} />
                <Skeleton variant="rect" width="100%" height={300} />
            </Box>
        )
    }

    let main = []
    for (let i in stInfo) {
        let data = stInfo[i]
        let s_title = data["Title"]
        let s_murl = data["MovieURL"]
        let s_purl = data["PictureURL"]
        let s_date = data["Date"]

        main.push(
            <LazyLoad
                key={i}
                height={200}
                offset={[-200, 0]}
                once
            >
                <Grid item xs={8} className={classes.stBody}>
                    <Card className={classes.stCard}>
                        <CardContent classes={{ root: classes.stCardContent }}>
                            <Typography component="p">
                                {s_date}
                            </Typography>
                            <Typography className={classes.stTitle} dangerouslySetInnerHTML={{ __html: s_title }} >
                            </Typography>
                            <Typography>
                                <img className={classes.stImg} src={s_purl} alt={"i" + i}></img>
                            </Typography>
                            <Typography component="p" className={classes.stDown}>
                                <Button
                                    variant="contained"
                                    color="primary"
                                    href={s_murl}
                                    target="_blank"
                                    classes={{
                                        root: classes.stDownCls
                                    }}
                                >
                                    Watch
                                </Button>
                            </Typography>
                        </CardContent>
                    </Card>
                </Grid>
            </LazyLoad>
        )
    }
    return (
        <Typography component="div">
            <AppBar className={classes.mainBar}>
                <Toolbar>
                    STchannel Movies
                </Toolbar>
            </AppBar>

            <Grid
                container
                spacing={2}
                direction="column"
                justify="center"
                alignItems="center"
                className={classes.mainGrid}
            >
                {loading ? <Loading /> : main}
                <Typography component="div" className={classes.stLoad} style={stLoadHidden}>
                    <Button
                        variant="contained"
                        classes={{
                            root: classes.stDownCls
                        }}
                        disabled={stLoadState}
                        onClick={() => getStList()}
                    >
                        {stLoadText}
                    </Button>
                </Typography>
            </Grid>

            <Snackbar
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'center'
                }}
                open={open}
                autoHideDuration={3000}
                onClose={() => barClose()}
                ContentProps={{
                    'aria-describedby': 'message-id',
                }}
                message={message}
            />
        </Typography>
    )
}
