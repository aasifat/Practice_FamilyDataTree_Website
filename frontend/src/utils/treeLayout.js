/**
 * Binary-Tree Layout Algorithm for Family Trees
 * - Left = Male lineage
 * - Right = Female lineage
 * - Couples treated as single unit
 * - Hierarchical positioning with consistent spacing
 */

/**
 * Calculate hierarchical generations
 * @param {Array} members - All people in tree
 * @returns {Object} Map of person ID to generation level
 */
export function calculateGenerations(members) {
  const generations = {};

  const assignGeneration = (personId, gen = 0) => {
    if (generations[personId] !== undefined) return generations[personId];

    generations[personId] = gen;

    const person = members.find((m) => m.id === personId);
    if (person?.father_id) {
      assignGeneration(person.father_id, gen - 1);
    }
    if (person?.mother_id) {
      assignGeneration(person.mother_id, gen - 1);
    }

    const children = members.filter(
      (m) => m.father_id === personId || m.mother_id === personId,
    );
    children.forEach((child) => assignGeneration(child.id, gen + 1));

    return gen;
  };

  members.forEach((member) => assignGeneration(member.id));
  return generations;
}

/**
 * Group members by generation
 * @param {Array} members
 * @param {Object} generations
 * @returns {Object} Members grouped by generation level
 */
export function groupByGeneration(members, generations) {
  const grouped = {};

  members.forEach((member) => {
    const gen = generations[member.id] || 0;
    if (!grouped[gen]) grouped[gen] = [];
    grouped[gen].push(member);
  });

  return grouped;
}

/**
 * Create couple units (married pairs)
 * @param {Array} members
 * @returns {Array} Array of couple units
 */
export function createCoupleUnits(members) {
  const couples = [];
  const processedIds = new Set();

  members.forEach((person) => {
    if (processedIds.has(person.id)) return;

    if (person.spouse_id) {
      const spouse = members.find((m) => m.id === person.spouse_id);
      if (spouse) {
        couples.push({
          person1: person,
          person2: spouse,
          id: `couple-${Math.min(person.id, spouse.id)}-${Math.max(
            person.id,
            spouse.id,
          )}`,
        });
        processedIds.add(person.id);
        processedIds.add(spouse.id);
      }
    }
  });

  return couples;
}

/**
 * Layout algorithm using D3-like hierarchical positioning
 * @param {Object} parameters
 * @returns {Object} Position map for nodes
 */
export function calculateLayout({
  members,
  generations,
  membersByGen,
  containerWidth = 1200,
  generationHeight = 150,
  siblingSpacing = 100,
}) {
  const positions = {};
  const sortedGens = Object.keys(membersByGen)
    .map(Number)
    .sort((a, b) => a - b);

  sortedGens.forEach((gen, genIndex) => {
    const membersInGen = membersByGen[gen] || [];
    const genY = genIndex * generationHeight + generationHeight / 2;

    // Calculate x positions for members in this generation
    const totalMembers = membersInGen.length;
    const totalWidth = (totalMembers - 1) * siblingSpacing + 100;
    const startX = Math.max(50, (containerWidth - totalWidth) / 2);

    membersInGen.forEach((member, index) => {
      const x = startX + index * siblingSpacing;
      positions[member.id] = {
        x,
        y: genY,
        generationLevel: gen,
      };
    });
  });

  return positions;
}

/**
 * Calculate binary positioning (left for male, right for female lineage)
 * @param {Object} parameters
 * @returns {Object} Position map with binary layout
 */
export function calculateBinaryLayout({
  members,
  generations,
  membersByGen,
  containerWidth = 1200,
  generationHeight = 150,
  lineageSpacing = 250,
}) {
  const positions = {};
  const sortedGens = Object.keys(membersByGen)
    .map(Number)
    .sort((a, b) => a - b);

  const centerX = containerWidth / 2;

  sortedGens.forEach((gen, genIndex) => {
    const membersInGen = membersByGen[gen] || [];
    const genY = genIndex * generationHeight + generationHeight / 2;

    // Separate by lineage
    const maleLineage = membersInGen.filter((m) => m.gender === "male");
    const femaleLineage = membersInGen.filter((m) => m.gender === "female");
    const otherLineage = membersInGen.filter(
      (m) => m.gender !== "male" && m.gender !== "female",
    );

    // Position male lineage on left
    maleLineage.forEach((member, index) => {
      const offset = index * 80 - (maleLineage.length * 80) / 2;
      positions[member.id] = {
        x: centerX - lineageSpacing + offset,
        y: genY,
        generationLevel: gen,
        lineage: "male",
      };
    });

    // Position female lineage on right
    femaleLineage.forEach((member, index) => {
      const offset = index * 80 - (femaleLineage.length * 80) / 2;
      positions[member.id] = {
        x: centerX + lineageSpacing + offset,
        y: genY,
        generationLevel: gen,
        lineage: "female",
      };
    });

    // Position other in middle
    otherLineage.forEach((member, index) => {
      const offset = index * 80 - (otherLineage.length * 80) / 2;
      positions[member.id] = {
        x: centerX + offset,
        y: genY,
        generationLevel: gen,
        lineage: "other",
      };
    });
  });

  return positions;
}

/**
 * Generate connection lines between parents and children
 * @param {Array} members
 * @param {Object} positions
 * @returns {Array} Array of line objects for rendering
 */
export function generateConnectionLines(members, positions) {
  const lines = [];

  members.forEach((person) => {
    const childPos = positions[person.id];
    if (!childPos) return;

    // Line from father
    if (person.father_id) {
      const fatherPos = positions[person.father_id];
      if (fatherPos) {
        const midY = (childPos.y + fatherPos.y) / 2;
        lines.push({
          type: "parent-to-child",
          startX: fatherPos.x,
          startY: fatherPos.y,
          midY,
          endX: childPos.x,
          endY: childPos.y,
          key: `line-father-${person.id}`,
        });
      }
    }

    // Line from mother
    if (person.mother_id) {
      const motherPos = positions[person.mother_id];
      if (motherPos) {
        const midY = (childPos.y + motherPos.y) / 2;
        lines.push({
          type: "parent-to-child",
          startX: motherPos.x,
          startY: motherPos.y,
          midY,
          endX: childPos.x,
          endY: childPos.y,
          key: `line-mother-${person.id}`,
        });
      }
    }

    // Line between spouses
    if (person.spouse_id && person.id < person.spouse_id) {
      // Only create line once per couple
      const spousePos = positions[person.spouse_id];
      if (spousePos) {
        lines.push({
          type: "spouse",
          startX: childPos.x,
          startY: childPos.y,
          endX: spousePos.x,
          endY: spousePos.y,
          key: `line-spouse-${person.id}-${person.spouse_id}`,
        });
      }
    }
  });

  return lines;
}

/**
 * Get root nodes (people without parents)
 * @param {Array} members
 * @returns {Array} Root members
 */
export function getRootNodes(members) {
  return members.filter((m) => !m.father_id && !m.mother_id);
}

/**
 * Get all descendants of a person
 * @param {Number} personId
 * @param {Array} members
 * @returns {Array} All descendants
 */
export function getDescendants(personId, members) {
  const descendants = [];
  const visited = new Set();

  const traverse = (id) => {
    if (visited.has(id)) return;
    visited.add(id);

    const children = members.filter(
      (m) => m.father_id === id || m.mother_id === id,
    );
    children.forEach((child) => {
      descendants.push(child);
      traverse(child.id);
    });
  };

  traverse(personId);
  return descendants;
}

/**
 * Validate tree integrity
 * @param {Array} members
 * @returns {Object} Validation result with errors and warnings
 */
export function validateTreeIntegrity(members) {
  const errors = [];
  const warnings = [];

  members.forEach((person) => {
    // Check: Child must have both parents
    if (
      (person.father_id || person.mother_id) &&
      !(person.father_id && person.mother_id)
    ) {
      errors.push(
        `Person "${person.name}" has only one parent. Both parents required.`,
      );
    }

    // Check: Parents must be spouses
    if (person.father_id && person.mother_id) {
      const father = members.find((m) => m.id === person.father_id);
      const mother = members.find((m) => m.id === person.mother_id);

      if (father && mother) {
        if (father.spouse_id !== mother.id && mother.spouse_id !== father.id) {
          warnings.push(
            `Person "${person.name}": Father and mother are not marked as spouses`,
          );
        }
      }
    }

    // Check: No duplicate relationships
    if (person.spouse_id && person.spouse_id === person.id) {
      errors.push(`Person "${person.name}" cannot be their own spouse`);
    }
  });

  return { valid: errors.length === 0, errors, warnings };
}
