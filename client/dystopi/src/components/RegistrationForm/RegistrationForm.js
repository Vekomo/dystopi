import React, {useState} from 'react';
function RegistrationForm(props) {
  const [state, setState] = useState({
    username : "",
    password: "",
    confirmPassword: ""
  })
  const handleChange = (e) => {
    const {id, value} = e.target
    setState(prevState => ({
      ...prevState,
      [id] : value
    }))
  }
  const handleSubmitClick = (e) => {
    e.preventDefault();
    if(state.password === state.confirmPassword) {
      sendDetailsToServer()
    } else {
      console.log('Passwords do not match');
    }
  }
  const sendDetailsToServer = () => {
    if(state.username.length && state.password.length) {
      const payload = {
        "username": state.username,
        "password": state.password,
      }
      //axios post here with error catch, but for now a console log.
      console.log('---> And then we\'d send the post request over.');
    } else {
      console.log('---> Something wrong with lengths of username/pass');
      console.log('Please enter valid username and password.');
    }
  }
  return(
        <div className="card col-12 col-lg-4 login-card mt-2 hv-center">
            <form>
                <div className="form-group text-left">
                <label htmlFor="exampleInputUsername">Username</label>
                <input type="username"
                       className="form-control"
                       id="username"
                       aria-describedby="usernameHelp"
                       placeholder="Enter username"
                       value={state.username}
                       onChange={handleChange}
                />
                </div>
                <div className="form-group text-left">
                    <label htmlFor="exampleInputPassword1">Password</label>
                    <input type="password"
                        className="form-control"
                        id="password"
                        placeholder="Password"
                        value={state.password}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group text-left">
                    <label htmlFor="exampleInputPassword1">Confirm Password</label>
                    <input type="password"
                        className="form-control"
                        id="confirmPassword"
                        placeholder="Confirm Password"
                        value={state.confirmPassword}
                        onChange={handleChange}
                    />
                </div>
                <button
                    type="submit"
                    className="btn btn-primary"
                    onClick={handleSubmitClick}
                >
                    Register
                </button>
            </form>
        </div>
    )
  }
export default RegistrationForm;
