import React from 'react'

import PropTypes from 'prop-types'
import { withStyles } from '@material-ui/styles'
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

import LinearProgress from '@material-ui/core/LinearProgress'
import SaveIcon from '@material-ui/icons/Save'
import CheckIcon from '@material-ui/icons/Check'
import SandIcon from '@material-ui/icons/HourglassEmpty'

const styles = ({
  mainBar: {
    position: 'fixed',
    backgroundColor: "#dd4b9a"
  },
  mainGrid: {
    position: 'relative',
    top: 100
  },
  stCard: {
    margin: 10,
    textAlign: "center",
  },
  stCardContent: {
    paddingBottom: 0
  },
  stDate: {
    padding: 10
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
    },
    '&.$Mui-disabled': {
      color: "#fff",
      backgroundColor: "#ab1e6b",
    }
  },
  stLoad: {
    padding: 10,
    height: 200,
  },
  progressRoot: {
    margin: "0 auto",
    maxWidth: 120
  },
  progressColor: {
    backgroundColor: "#f2c2dc",
  },
  progressBarColor: {
    backgroundColor: "#dd4b9a"
  },
  SkeletonCls: {
    margin: '10px auto'
  }
})

class Index extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      stinfo: [],
      minid: 9999,
      message: "",
      open: false,
      loading: true,
      vertical: 'top',
      horizontal: 'center',
      stLoadHidden: {
        display: "none"
      },
      progDisplay: false,
      stDownMsg: "download",
      stLoadState: false,
      stLoadText: "more",
      stBtnIcon: <SaveIcon />
    }
  }

  componentDidMount() {
    this.getStList()
  }
  getStList() {
    this.setState({
      stLoadState: true,
    })
    let url = ""
    if (this.state.minid === 9999) {
      url = `http://localhost:61800/stlist`
    } else {
      url = `http://localhost:61800/stlist?next_id=${this.state.minid}`
    }
    fetch(
      url,
      {
        method: "GET",
        dataType: "json",
      }).then(res => res.json())
      .then(response => {
        if (response.status === 1) {
          let sids = []
          let data = this.state.stinfo
          let add_data = response.data
          for (let i in add_data) {
            let d = add_data[i]
            sids.push(parseInt(d["ID"]))
          }
          let minid = Math.min.apply(null, sids)
          let newdata = data.concat(add_data)
          this.setState({
            stinfo: newdata,
            minid: minid,
            stLoadHidden: {
              display: "block"
            },
            loading: false,
            stLoadState: false,
          })
        } else {
          this.setState({
            stLoadHidden: {
              display: "block"
            },
            stLoadText: "no data",
            loading: false,
            stLoadState: true,
          })
        }
      }).catch(
        () => this.setState({
          open: true,
          message: "Server Error",
          loading: false,
          stLoadState: false,
        })
      )
  }

  getMedia(murl, btnid) {
    let mData = {
      "url": murl
    }
    // www2.uliza.jp 
    if (murl.indexOf("www2.uliza.jp") === -1) {
      this.setState({
        open: true,
        message: "地址为外链 " + murl,
      })
      return false
    }
    let downState = 'stDownState' + btnid
    let downMsg = 'stDownMsg' + btnid
    let progDisplay = 'progDisplay' + btnid
    this.setState({
      [downState]: true,
      [downMsg]: "waiting",
      stBtnIcon: <SandIcon />
    })
    this.fakeProgress(btnid)
    fetch(
      `http://localhost:61800/stdown`,
      {
        method: 'POST',
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(mData),
        dataType: 'json',
      }).then(res => res.json())
      .then(data => {
        if (data.status === 1) {
          this.setState({
            [downState]: true,
            [downMsg]: "finished",
            [progDisplay]: false,
            stBtnIcon: <CheckIcon />
          })
        } else {
          this.setState({
            [downState]: false,
            [downMsg]: "download",
            open: true,
            message: "Error",
            [progDisplay]: false,
            stBtnIcon: <SaveIcon />
          })
        }
      }).catch(
        () => this.setState({
          [downState]: false,
          [downMsg]: "download",
          open: true,
          message: "Server Error",
          [progDisplay]: false,
          stBtnIcon: <SaveIcon />
        })
      )
  }

  fakeProgress(i) {
    let linear = "linear" + i
    let progDisplay = "progDisplay" + i
    if (this.state[linear] === undefined || this.state[linear] !== 0) {
      this.setState({
        [progDisplay]: true,
        [linear]: 0
      })
    }
    
    const timer = setInterval(() => {
      // console.log(linear, progDisplay, this.state[linear], this.state[progDisplay])
        if (this.state[linear] === 100) {
          clearInterval(timer)
        }
        this.setState({
          [linear]: Math.min(this.state[linear] + Math.random() * 2, 100),
        })
      }, 1000)
    return timer
    
  }

  handleClose(event, reason) {
    if (reason === 'clickaway') {
      return
    }
    this.setState({
      open: false
    })
  }

  componentWillUnmount() {
    clearInterval(this.fakeProgress)
  }

  render() {
    const { classes } = this.props
    const stinfo = this.state.stinfo
    let main = []

    function Loading() {
      return (
        <Box width="65%" >
          <Skeleton width="24%" height={30} className={classes.SkeletonCls} />
          <Skeleton width="100%" height={10} />
          <Skeleton width="100%" height={10} />
          <Skeleton width="80%" height={10} />
          <Skeleton variant="rect" width="100%" height={400} />
          <Skeleton width="14%" height={46} className={classes.SkeletonCls} />
        </Box>
      )
    }

    for (let i in stinfo) {
      let data = stinfo[i]
      let s_title = data["Title"]
      let s_murl = data["MovieURL"]
      let s_purl = data["PictureURL"]
      let s_date = data["Date"]

      let downMsg = this.state[`stDownMsg${i}`]
      if (downMsg === undefined) {
        downMsg = "download"
      }

      main.push(
        <LazyLoad
          key={i}
          height={200}
          offset={[-200, 0]}
          once
          placeholder={<Loading />}
        >
          <Grid item xs={8} >
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
                <Typography component="div" style={this.state[`progDisplay${i}`] ? {display: "block"} : {display: "none"}} >
                  <LinearProgress
                    variant="determinate"
                    classes={{
                      root: classes.progressRoot,
                      colorPrimary: classes.progressColor,
                      barColorPrimary: classes.progressBarColor
                    }}
                    value={this.state[`linear${i}`] ? this.state[`linear${i}`] : 0}
                  />
                </Typography>
                <Typography component="p" className={classes.stDown}>
                  <Button
                    disabled={this.state[`stDownState${i}`]}
                    variant="contained"
                    color="primary"
                    onClick={this.getMedia.bind(this, s_murl, i)}
                    classes={{
                      root: classes.stDownCls
                    }}
                    startIcon={this.state.stBtnIcon}
                  >
                    {downMsg}
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
          {this.state.loading ? <Loading /> : main}
          <Typography component="div" className={classes.stLoad} style={this.state.stLoadHidden}>
            <Button
              variant="contained"
              classes={{
                root: classes.stDownCls
              }}
              disabled={this.state.stLoadState}
              onClick={this.getStList.bind(this)}
            >
              {this.state.stLoadText}
            </Button>
          </Typography>
        </Grid>

        <Snackbar
          anchorOrigin={{
            vertical: this.state.vertical,
            horizontal: this.state.horizontal
          }}
          open={this.state.open}
          autoHideDuration={3000}
          onClose={this.handleClose.bind(this)}
          ContentProps={{
            'aria-describedby': 'message-id',
          }}
          message={this.state.message}
        />
      </Typography>
    )
  }
}

Index.propTypes = {
  classes: PropTypes.object.isRequired,
}

export default withStyles(styles)(Index)
