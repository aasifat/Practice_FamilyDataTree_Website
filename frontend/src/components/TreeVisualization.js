import React, { useMemo, useRef, useState, useEffect } from "react";
import { calculateGenerations, groupByGeneration } from "../utils/treeLayout";
import "./TreeVisualization.css";

const TreeVisualization = ({ tree, onDeleteMember }) => {
  const containerRef = useRef(null);
  const cardRefsRef = useRef({});
  const [cardPositions, setCardPositions] = useState({});
  const [svgHeight, setSvgHeight] = useState(0);
  const [zoom, setZoom] = useState(1);
  const [pan, setPan] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });

  // Calculate generations
  const generations = useMemo(() => {
    const members = tree.members || [];
    return calculateGenerations(members);
  }, [tree.members]);

  // Group by generation
  const membersByGen = useMemo(() => {
    const members = tree.members || [];
    return groupByGeneration(members, generations);
  }, [tree.members, generations]);

  // Get sorted generations
  const sortedGens = useMemo(() => {
    return Object.keys(membersByGen)
      .map(Number)
      .sort((a, b) => a - b);
  }, [membersByGen]);

  // Update positions from DOM
  const updatePositions = () => {
    if (!containerRef.current) return;

    const positions = {};
    const containerRect = containerRef.current.getBoundingClientRect();

    Object.keys(cardRefsRef.current).forEach((memberId) => {
      const el = cardRefsRef.current[memberId];
      if (!el) return;

      const rect = el.getBoundingClientRect();
      positions[memberId] = {
        x: rect.left - containerRect.left + rect.width / 2,
        y: rect.top - containerRect.top + rect.height / 2,
      };
    });

    setCardPositions(positions);
    setSvgHeight(containerRef.current.scrollHeight);
  };

  // Update positions when layout changes
  useEffect(() => {
    const timeout = setTimeout(updatePositions, 50);
    return () => clearTimeout(timeout);
  }, [membersByGen]);

  // Resize observer
  useEffect(() => {
    if (!containerRef.current) return;

    const observer = new ResizeObserver(() => {
      updatePositions();
    });

    observer.observe(containerRef.current);
    return () => observer.disconnect();
  }, []);

  // Connection lines
  const connectionLines = useMemo(() => {
    const lines = [];
    const members = tree.members || [];

    members.forEach((person) => {
      if (!person.father_id && !person.mother_id) return;

      const childPos = cardPositions[person.id];

      // Line from father
      if (person.father_id) {
        const fatherPos = cardPositions[person.father_id];
        if (childPos && fatherPos) {
          const midY = (childPos.y + fatherPos.y) / 2;
          const midX = (fatherPos.x + childPos.x) / 2;
          const relationLabel = person.gender === "male" ? "Son" : "Daughter";
          lines.push(
            <g key={`line-father-${person.id}`}>
              <line
                x1={fatherPos.x}
                y1={fatherPos.y}
                x2={fatherPos.x}
                y2={midY}
                stroke="#666"
                strokeWidth="2"
              />
              <line
                x1={fatherPos.x}
                y1={midY}
                x2={childPos.x}
                y2={midY}
                stroke="#666"
                strokeWidth="2"
              />
              <line
                x1={childPos.x}
                y1={midY}
                x2={childPos.x}
                y2={childPos.y}
                stroke="#666"
                strokeWidth="2"
              />
              <text
                x={midX}
                y={midY - 5}
                textAnchor="middle"
                fontSize="11"
                fill="#333"
                fontWeight="bold"
                backgroundColor="white"
                style={{ pointerEvents: "none" }}
              >
                {relationLabel}
              </text>
            </g>,
          );
        }
      }

      // Line from mother
      if (person.mother_id) {
        const motherPos = cardPositions[person.mother_id];
        if (childPos && motherPos) {
          const midY = (childPos.y + motherPos.y) / 2;
          const midX = (motherPos.x + childPos.x) / 2;
          const relationLabel = person.gender === "male" ? "Son" : "Daughter";
          lines.push(
            <g key={`line-mother-${person.id}`}>
              <line
                x1={motherPos.x}
                y1={motherPos.y}
                x2={motherPos.x}
                y2={midY}
                stroke="#888"
                strokeWidth="2"
                strokeDasharray="5,5"
              />
              <line
                x1={motherPos.x}
                y1={midY}
                x2={childPos.x}
                y2={midY}
                stroke="#888"
                strokeWidth="2"
                strokeDasharray="5,5"
              />
              <line
                x1={childPos.x}
                y1={midY}
                x2={childPos.x}
                y2={childPos.y}
                stroke="#888"
                strokeWidth="2"
                strokeDasharray="5,5"
              />
              <text
                x={midX}
                y={midY - 5}
                textAnchor="middle"
                fontSize="11"
                fill="#555"
                fontWeight="bold"
                style={{ pointerEvents: "none" }}
              >
                {relationLabel}
              </text>
            </g>,
          );
        }
      }
    });

    // Spouse connection lines
    members.forEach((person) => {
      if (person.spouse_id && person.id < person.spouse_id) {
        const pos1 = cardPositions[person.id];
        const pos2 = cardPositions[person.spouse_id];

        if (pos1 && pos2) {
          const midX = (pos1.x + pos2.x) / 2;
          const midY = (pos1.y + pos2.y) / 2;
          lines.push(
            <g key={`spouse-${person.id}-${person.spouse_id}`}>
              <line
                x1={pos1.x}
                y1={pos1.y}
                x2={pos2.x}
                y2={pos2.y}
                stroke="#e74c3c"
                strokeWidth="3"
              />
              <text
                x={midX}
                y={midY - 8}
                textAnchor="middle"
                fontSize="10"
                fill="#e74c3c"
                fontWeight="bold"
                style={{ pointerEvents: "none" }}
              >
                Spouse
              </text>
            </g>,
          );
        }
      }
    });

    return lines;
  }, [cardPositions, tree.members]);

  // Render card
  const renderCard = (person) => {
    const initials = person.name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase();

    const genderColor = person.gender === "male" ? "#3498db" : "#e91e63";

    return (
      <div
        key={person.id}
        className="tree-card"
        style={{ borderColor: genderColor }}
        ref={(el) => {
          if (el) cardRefsRef.current[person.id] = el;
        }}
      >
        <div className="card-initials" style={{ backgroundColor: genderColor }}>
          {initials}
        </div>
        <div className="card-content">
          <h3 className="card-name">{person.name}</h3>
          <p className="card-gender">{person.gender}</p>
          {person.spouse_id && <p className="card-spouse">♥ Married</p>}
          <button
            className="delete-btn"
            onClick={() => onDeleteMember(person.id)}
            title="Delete member"
          >
            ✕
          </button>
        </div>
      </div>
    );
  };

  // Pan and zoom handlers
  const handleWheel = (e) => {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault();
      const delta = e.deltaY > 0 ? 0.9 : 1.1;
      setZoom((z) => Math.max(0.5, Math.min(3, z * delta)));
    }
  };

  const handleMouseDown = (e) => {
    if (e.button === 2 || (e.ctrlKey && !e.ctrlKey)) {
      setIsDragging(true);
      setDragStart({ x: e.clientX - pan.x, y: e.clientY - pan.y });
    }
  };

  const handleMouseMove = (e) => {
    if (isDragging) {
      setPan({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  // Empty state
  if (!tree.members || tree.members.length === 0) {
    return (
      <div className="tree-container empty">
        <p>No members in the tree yet. Add some members to visualize!</p>
      </div>
    );
  }

  return (
    <div
      className="tree-container"
      ref={containerRef}
      onWheel={handleWheel}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
      style={{ cursor: isDragging ? "grabbing" : "grab" }}
    >
      <div className="tree-controls">
        <button
          onClick={() => setZoom(1)}
          className="control-btn"
          title="Reset zoom"
        >
          🔍 Reset
        </button>
        <span className="zoom-level">{Math.round(zoom * 100)}%</span>
      </div>

      <div className="tree-title">Family Tree Structure (Binary Layout)</div>

      <div className="tree-wrapper">
        <svg
          className="tree-connections"
          width="100%"
          height={svgHeight || 800}
          style={{
            position: "absolute",
            top: 0,
            left: 0,
            pointerEvents: "none",
            transform: `scale(${zoom}) translate(${pan.x}px, ${pan.y}px)`,
            transformOrigin: "0 0",
          }}
        >
          {connectionLines}
        </svg>

        <div
          className="generations-container"
          style={{
            transform: `scale(${zoom}) translate(${pan.x}px, ${pan.y}px)`,
            transformOrigin: "0 0",
          }}
        >
          {sortedGens.map((gen) => (
            <div key={gen} className="generation">
              <div className="gen-label">Gen {gen + 1}</div>

              <div className="members-row">
                {membersByGen[gen]?.map(renderCard)}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="legend">
        <div className="legend-item">
          <span className="color-box" style={{ backgroundColor: "#3498db" }}>
            ▯
          </span>
          Male
        </div>
        <div className="legend-item">
          <span className="color-box" style={{ backgroundColor: "#e91e63" }}>
            ▯
          </span>
          Female
        </div>
        <div className="legend-item">
          <svg width="20" height="20" style={{ display: "inline" }}>
            <line
              x1="0"
              y1="10"
              x2="20"
              y2="10"
              stroke="#666"
              strokeWidth="2"
            />
          </svg>
          Father-Child
        </div>
        <div className="legend-item">
          <svg width="20" height="20" style={{ display: "inline" }}>
            <line
              x1="0"
              y1="10"
              x2="20"
              y2="10"
              stroke="#888"
              strokeWidth="2"
              strokeDasharray="3,3"
            />
          </svg>
          Mother-Child
        </div>
        <div className="legend-item">
          <svg width="20" height="20" style={{ display: "inline" }}>
            <line
              x1="0"
              y1="10"
              x2="20"
              y2="10"
              stroke="#e74c3c"
              strokeWidth="3"
            />
          </svg>
          Spouse
        </div>
      </div>
    </div>
  );
};

export default TreeVisualization;
