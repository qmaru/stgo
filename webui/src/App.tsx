import React, { useState, useEffect, useCallback } from 'react'
import { makeStyles } from '@mui/styles'
import { createTheme, ThemeProvider } from '@mui/material/styles'

import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import Typography from '@mui/material/Typography'
import Grid from '@mui/material/Grid'
import Button from '@mui/material/Button'
import Skeleton from '@mui/lab/Skeleton'
import LazyLoad from 'react-lazyload'
import Box from '@mui/material/Box'
import AppBar from '@mui/material/AppBar'
import Toolbar from '@mui/material/Toolbar'

const theme = createTheme({
    palette: {
        primary: {
            main: "#dd4b9a"
        },
    },
    components: {
        MuiGrid: {
            styleOverrides: {
                root: {
                    maxWidth: "100%"
                }
            }
        }
    }
});

const useStyles = makeStyles({
    mainWrapper: {
        position: "relative",
        top: 80
    },
    stBody: {
        padding: 8
    },
    stCard: {
        textAlign: "center"
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
    stLoad: {
        padding: 10,
        height: 50,
    },
    SkeletonCls: {
        margin: '10px auto'
    }
})


export default function Index() {
    const classes = useStyles()
    // STChannel 数据
    const [STMovieData, setSTMovieData] = useState<any>([])
    const [prevID, setPrevID] = useState<number>(0)
    // UI 数据
    const [loading, setLoading] = useState<boolean>(true)
    const [STMoreDisalbed, setSTMoreDisalbed] = useState<boolean>(false)
    const [STMoreShow, setSTMoreShow] = useState<any>({ display: "none" })
    const [STMoreText, setSTMoreText] = useState<string>("more")

    const GetSTData = () => {
        var url: string = `http://127.0.0.1:61800/stlist?next_id=${prevID}`
        fetch(url, { method: "GET" })
            .then(res => res.json())
            .then(response => {
                if (response.status === 1) {
                    let STData: any = response.data
                    GetEarliestID(STData)
                    let newData: any = STMovieData.concat(STData)
                    setSTMovieData(newData)
                } else {
                    setLoading(false)
                    setSTMoreShow({ display: "block" })
                    setSTMoreDisalbed(true)
                    setSTMoreText("no data")
                }
            }).catch(
                () => {
                    setLoading(false)
                    setSTMoreShow({ display: "block" })
                    setSTMoreDisalbed(true)
                    setSTMoreText("Server Error")
                }
            )
    }

    // 获取最前的 ID
    const GetEarliestID = (STData: any) => {
        let STIDList: any = []
        for (let i in STData) {
            let sData: any = STData[i]
            STIDList.push(parseInt(sData["ID"]))
        }
        let earliestID: number = Math.min.apply(null, STIDList)
        setPrevID(earliestID)
    }

    const LoadFirstData = useCallback(() => {
        const url: string = `http://127.0.0.1:61800/stlist`
        fetch(url, { method: "GET" })
            .then(res => res.json())
            .then(response => {
                if (response.status === 1) {
                    let STData: any = response.data
                    GetEarliestID(STData)
                    setSTMovieData(STData)
                    setSTMoreShow({ display: "block" })
                    setLoading(false)
                } else {
                    setLoading(false)
                    setSTMoreShow({ display: "block" })
                    setSTMoreDisalbed(true)
                    setSTMoreText("no data")
                }
            }).catch(
                () => {
                    setLoading(false)
                    setSTMoreShow({ display: "block" })
                    setSTMoreDisalbed(true)
                    setSTMoreText("Server Error")
                }
            )
    }, [])

    useEffect(() => {
        LoadFirstData()
    }, [LoadFirstData])

    const Loading = () => {
        return (
            <Box width="100%" >
                <Skeleton animation="wave" width="24%" height={60} className={classes.SkeletonCls} />
                <Skeleton animation="wave" width="100%" height={30} />
                <Skeleton animation="wave" width="100%" height={30} />
                <Skeleton animation="wave" width="80%" height={30} />
                <Skeleton animation="wave" variant="rectangular" width="100%" height={300} />
            </Box>
        )
    }

    let main = []
    for (let i in STMovieData) {
        let data = STMovieData[i]
        let s_title = data["Title"]
        let s_murl = data["MovieURL"]
        let s_purl = data["PictureURL"]
        let s_date = data["Date"]

        main.push(
            <LazyLoad
                key={"stcard" + i}
                height={200}
                offset={[-200, 0]}
                once
            >
                <Grid item xs={8} className={classes.stBody}>
                    <Card className={classes.stCard}>
                        <CardContent>
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
        <ThemeProvider theme={theme}>
            <div>
                <AppBar position="fixed">
                    <Toolbar>
                        STchannel Movies
                    </Toolbar>
                </AppBar>

                <div className={classes.mainWrapper}>
                    <Grid
                        container
                        direction="column"
                        justifyContent="center"
                        alignItems="center"
                    >
                        {loading ? <Loading /> : main}
                        <div className={classes.stLoad} style={STMoreShow}>
                            <Button
                                variant="contained"
                                disabled={STMoreDisalbed}
                                onClick={() => GetSTData()}
                            >
                                {STMoreText}
                            </Button>
                        </div>
                    </Grid>
                </div>
            </div>
        </ThemeProvider>
    )
}
