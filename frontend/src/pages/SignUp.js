import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  requestOTPForSignup,
  verifyOTPAndSignup,
} from "../services/authService";
import { useAuth } from "../context/AuthContext";
import {
  MdEmail,
  MdLock,
  MdPerson,
  MdVisibilityOff,
  MdVisibility,
} from "react-icons/md";
import "../styles/Auth.css";

const SignUp = () => {
  const navigate = useNavigate();
  const { setUser } = useAuth();
  const [step, setStep] = useState("form");
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
  });
  const [otp, setOtp] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleOTPChange = (e) => {
    setOtp(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      await requestOTPForSignup({
        email: formData.email,
        password: formData.password,
        first_name: formData.firstName,
        last_name: formData.lastName,
      });
      setStep("otp");
    } catch (err) {
      setError(err.response?.data?.error || "Failed to request OTP");
    } finally {
      setLoading(false);
    }
  };

  const handleOTPSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const result = await verifyOTPAndSignup({
        email: formData.email,
        otp: otp,
        password: formData.password,
        first_name: formData.firstName,
        last_name: formData.lastName,
      });
      setUser(result.user);
      navigate("/");
    } catch (err) {
      setError(err.response?.data?.error || "OTP verification failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-background">
        <div className="auth-gradient-1"></div>
        <div className="auth-gradient-2"></div>
      </div>

      <div className="auth-content">
        <div className="auth-card">
          <div className="auth-header">
            <h1>{step === "form" ? "Create Account" : "Verify Email"}</h1>
            <p>
              {step === "form"
                ? "Join us and explore your family tree"
                : "Enter the code sent to your email"}
            </p>
          </div>

          {error && (
            <div className="alert alert-error">
              <span>{error}</span>
            </div>
          )}

          {step === "form" ? (
            <form onSubmit={handleSubmit} className="auth-form">
              <div className="form-group">
                <label className="form-label">First Name</label>
                <div className="input-wrapper">
                  <MdPerson className="input-icon" size={18} />
                  <input
                    type="text"
                    name="firstName"
                    value={formData.firstName}
                    onChange={handleChange}
                    placeholder="John"
                    required
                    className="form-input"
                  />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Last Name</label>
                <div className="input-wrapper">
                  <MdPerson className="input-icon" size={18} />
                  <input
                    type="text"
                    name="lastName"
                    value={formData.lastName}
                    onChange={handleChange}
                    placeholder="Doe"
                    required
                    className="form-input"
                  />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Email Address</label>
                <div className="input-wrapper">
                  <MdEmail className="input-icon" size={18} />
                  <input
                    type="email"
                    name="email"
                    value={formData.email}
                    onChange={handleChange}
                    placeholder="you@example.com"
                    required
                    className="form-input"
                  />
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Password</label>
                <div className="input-wrapper">
                  <MdLock className="input-icon" size={18} />
                  <input
                    type={showPassword ? "text" : "password"}
                    name="password"
                    value={formData.password}
                    onChange={handleChange}
                    placeholder="••••••••"
                    required
                    className="form-input"
                  />
                  <button
                    type="button"
                    className="password-toggle"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <MdVisibilityOff size={18} />
                    ) : (
                      <MdVisibility size={18} />
                    )}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="btn btn-primary btn-block"
              >
                {loading ? (
                  <>
                    <span className="loader"></span>
                    Sending OTP...
                  </>
                ) : (
                  "Continue"
                )}
              </button>
            </form>
          ) : (
            <form onSubmit={handleOTPSubmit} className="auth-form">
              <div
                className="alert alert-info"
                style={{ marginBottom: "1.5rem" }}
              >
                <span>
                  An OTP has been sent to <strong>{formData.email}</strong>
                </span>
              </div>

              <div className="form-group">
                <label className="form-label">Enter OTP Code</label>
                <input
                  type="text"
                  value={otp}
                  onChange={handleOTPChange}
                  placeholder="000000"
                  maxLength="6"
                  required
                  className="form-input"
                  style={{ textAlign: "center", letterSpacing: "0.5rem" }}
                />
              </div>

              <button
                type="submit"
                disabled={loading}
                className="btn btn-primary btn-block"
              >
                {loading ? (
                  <>
                    <span className="loader"></span>
                    Verifying...
                  </>
                ) : (
                  "Verify & Sign Up"
                )}
              </button>

              <button
                type="button"
                onClick={() => {
                  setStep("form");
                  setOtp("");
                  setError("");
                }}
                className="btn btn-outline btn-block"
              >
                Back
              </button>
            </form>
          )}

          <div className="auth-divider">
            <span>{step === "form" ? "Already have an account?" : ""}</span>
          </div>

          {step === "form" && (
            <Link to="/login" className="btn btn-outline btn-block">
              Sign In Instead
            </Link>
          )}
        </div>

        <div className="auth-info">
          <div className="info-card">
            <div className="info-icon">🌳</div>
            <h3>Family Tree Manager</h3>
            <p>Connect and visualize your family history in one place</p>
          </div>
          <div className="info-card">
            <div className="info-icon">✨</div>
            <h3>Easy to Use</h3>
            <p>Build your family tree with just a few clicks</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SignUp;
