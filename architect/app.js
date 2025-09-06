// Learning Architect - Complete JavaScript Application Logic (Fixed with AI Prompts)

// Application data (comprehensive dataset)
const appData = {
    "dashboard_stats": {
        "overall_progress": 88,
        "weekly_progress": 85,
        "daily_streak": 12,
        "total_study_time": 1247,
        "concepts_mastered": 89,
        "quiz_accuracy": 76,
        "retention_rate": 82
    },
    "todays_study_plan": [
        {
            "type": "weak_area_focus",
            "title": "Calculus Integration",
            "duration": "15 min",
            "difficulty": "medium",
            "action": "Start"
        },
        {
            "type": "spaced_repetition",
            "title": "Physics Formulas",
            "cards": "23 cards due",
            "action": "Review"
        },
        {
            "type": "problem_gymnasium",
            "title": "Advanced challenges",
            "remaining": "3 remaining",
            "action": "Practice"
        }
    ],
    "spaced_repetition_due": [
        {
            "topic": "Physics Formulas",
            "cards": 15,
            "difficulty": "Medium",
            "due": "Due Today"
        },
        {
            "topic": "French Verbs",
            "cards": 8,
            "difficulty": "Easy",
            "due": "Due Today"
        },
        {
            "topic": "Programming Concepts",
            "cards": 12,
            "difficulty": "Hard",
            "due": "Tomorrow"
        }
    ],
    "weak_areas": [
        {
            "topic": "Calculus Integration",
            "confidence": 45,
            "accuracy": 52,
            "priority": "High"
        },
        {
            "topic": "Historical Dates",
            "confidence": 38,
            "accuracy": 41,
            "priority": "High"
        },
        {
            "topic": "Spanish Vocabulary",
            "confidence": 62,
            "accuracy": 58,
            "priority": "Medium"
        }
    ],
    "weekly_analytics": {
        "days": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
        "learning_minutes": [25, 45, 60, 30, 55, 35, 50],
        "retention_scores": [7.4, 8.0, 8.6, 7.6, 8.4, 7.8, 8.2]
    },
    "learning_pathways": [
        {
            "name": "Foundation Through Discovery",
            "description": "Build intuitive understanding through guided exploration",
            "icon": "🔍",
            "progress": 65,
            "examples": [
                {"topic": "Calculus", "problem": "Why do we need derivatives?", "solution": "To find rates of change in real-world situations"},
                {"topic": "Chemistry", "problem": "What makes reactions happen?", "solution": "Energy changes and molecular collisions"}
            ]
        },
        {
            "name": "Applied Mastery Laboratory",
            "description": "Project-based progression with real-world applications",
            "icon": "🧪",
            "progress": 42,
            "features": ["Step-by-step guided practice", "70%+ success rate targeting"]
        },
        {
            "name": "Problem-Solving Gymnasium",
            "description": "Challenge levels from Novice to Expert mastery",
            "icon": "💪",
            "progress": 28,
            "difficulties": ["Novice", "Intermediate", "Advanced", "Expert"],
            "classic_problems": [
                {"problem": "Find area of circle with radius 5", "solution": "π × 5² = 25π ≈ 78.54", "difficulty": "Novice"},
                {"problem": "Solve 2x + 3 = 11", "solution": "x = 4", "difficulty": "Novice"}
            ]
        },
        {
            "name": "Mistake-Based Learning Lab",
            "description": "Error analysis and correction with anti-patterns",
            "icon": "🔧",
            "progress": 15,
            "misconceptions": [
                {"topic": "Math", "error": "(a+b)² = a² + b²", "correction": "(a+b)² = a² + 2ab + b²", "explanation": "Missing cross term 2ab"},
                {"topic": "Learning", "error": "Rereading is effective", "correction": "Active recall is much better", "explanation": "Familiarity ≠ knowledge"}
            ]
        },
        {
            "name": "Integration & Transfer System",
            "description": "Long-term retention with cross-domain connections",
            "icon": "🔗",
            "progress": 35,
            "methods": ["Spaced repetition scheduling", "Transfer learning optimization"]
        }
    ],
    "memory_techniques": {
        "memory_palaces": [
            {
                "name": "My House",
                "rooms": ["Front Door", "Living Room", "Kitchen", "Bedroom", "Bathroom"],
                "capacity": 20,
                "items": [
                    {"location": "Front Door", "item": "Milk", "image": "Giant milk bottle greeting visitors"},
                    {"location": "Living Room", "item": "Bread", "image": "Loaf of bread lounging on sofa"},
                    {"location": "Kitchen", "item": "Eggs", "image": "Eggs cooking themselves on stove"}
                ]
            },
            {
                "name": "School Campus",
                "rooms": ["Main Gate", "Library", "Cafeteria", "Gym", "Classroom"],
                "capacity": 25,
                "items": [
                    {"location": "Main Gate", "item": "Napoleon 1769", "image": "Napoleon on horseback at entrance"},
                    {"location": "Library", "item": "WWI 1914", "image": "Soldiers studying war books"},
                    {"location": "Cafeteria", "item": "Moon Landing 1969", "image": "Astronauts eating space food"}
                ]
            }
        ],
        "peg_system": {
            "pegs": {
                "1": {"word": "Tie", "image": "Necktie"},
                "2": {"word": "Noah", "image": "Man with long beard"},
                "3": {"word": "Ma", "image": "Mother figure"},
                "4": {"word": "Rye", "image": "Bottle of rye whiskey"},
                "5": {"word": "Law", "image": "Police badge"},
                "6": {"word": "Shoe", "image": "Leather shoe"},
                "7": {"word": "Cow", "image": "Black and white cow"},
                "8": {"word": "Ivy", "image": "Ivy growing on wall"},
                "9": {"word": "Bee", "image": "Buzzing honeybee"},
                "10": {"word": "Toes", "image": "Ten toes on feet"}
            },
            "sample_lists": [
                {
                    "name": "Shopping List",
                    "items": ["Milk", "Bread", "Apples", "Cheese", "Butter"],
                    "associations": [
                        "Tie soaked in milk",
                        "Noah carrying giant bread loaf",
                        "Ma juggling red apples",
                        "Rye bottle made of cheese",
                        "Police badge buttered up"
                    ]
                }
            ]
        },
        "link_method": {
            "sample_chains": [
                {
                    "topic": "Solar System Planets",
                    "chain": [
                        "Sun burning brightly",
                        "Mercury messenger running fast from sun",
                        "Venus goddess of beauty admiring Mercury",
                        "Earth our home planet greeting Venus",
                        "Mars red warrior threatening Earth"
                    ]
                }
            ]
        },
        "substitute_words": {
            "examples": [
                {"abstract": "Economics", "substitute": "Egg-onomics", "image": "Eggs managing money"},
                {"abstract": "Philosophy", "substitute": "Fill-loss-of-fee", "image": "Filing cabinet losing fees"},
                {"abstract": "Psychology", "substitute": "Sigh-college", "image": "Student sighing at college"},
                {"abstract": "Democracy", "substitute": "Dome-o-crazy", "image": "Crazy dome building"}
            ]
        },
        "chunking_exercises": {
            "description": "Pattern recognition and information grouping training",
            "examples": [
                {"pattern": "Phone Numbers", "chunked": "555-123-4567 → 555 | 123 | 4567"},
                {"pattern": "Credit Cards", "chunked": "1234567890123456 → 1234 | 5678 | 9012 | 3456"}
            ]
        }
    },
    "error_patterns": [
        {"topic": "Math Errors", "description": "Tendency to forget cross-terms in binomial expansion", "frequency": "High"},
        {"topic": "Physics Misconceptions", "description": "Confusing weight with mass in calculations", "frequency": "Medium"},
        {"topic": "Learning Strategy", "description": "Over-reliance on passive reading vs active recall", "frequency": "High"}
    ],
    "learning_strategies": [
        {"name": "Active Recall", "description": "Testing yourself significantly improves retention"},
        {"name": "Spaced Repetition", "description": "Optimal timing prevents forgetting curve"},
        {"name": "Memory Palace", "description": "Visual-spatial memory enhancement"}
    ]
};

// AI Prompts for Learning Pathways - Complete Implementation
const aiPrompts = {
    "foundation_discovery": `# 🔍 Foundation Through Discovery - AI Learning Coach

You are an **Adaptive Mastery Learning System (AMLS)** designed to build intuitive understanding through guided exploration with multi-dimensional assessment.

## Core Identity
You are an **Adaptive Learning Architect** combining cognitive science, emotional intelligence, and real-time personalization to create optimal learning experiences for topic: **{TOPIC}**

## Your Approach
- **Multi-dimensional assessment**: Assess cognitive load, emotional state, prior knowledge, and learning style
- **Real-time adaptation**: Adjust difficulty based on working memory capacity and emotional state
- **Cognitive load management**: Maximum 4 new elements per session
- **Progressive complexity**: 15-20% difficulty increase per level

## Session Protocol
\`\`\`
Initialize:
cognitive_load = assess_working_memory()
emotional_state = detect_affect_and_motivation() 
prior_knowledge = diagnostic_probe()

Adapt:
if cognitive_load > 0.7: reduce_complexity()
if emotional_state == "frustrated": switch_to_scaffolding()
if mastery_level > 0.85: increase_difficulty_and_interleave()
\`\`\`

## Learning Sequence
1. **Pattern Recognition** - Use analogies to build mental models
2. **Intuitive Building** - Connect to existing knowledge with visual mapping
3. **Active Application** - Immediate practice with spaced intervals
4. **Retrieval Practice** - 85% difficulty sweet spot targeting

## Spaced Repetition Schedule
- **20min → 1d → 3d → 1w → 2w → 1m → 3m**
- **Interleaving**: Mix 3-4 related concepts
- **Elaborative Rehearsal**: Connect to existing knowledge networks

## Adaptation Rules
- Accuracy < 50% → Simplify and scaffold
- Accuracy > 85% & Fast → Increase complexity  
- Plateau detected → Change modality/approach
- High engagement → Extend session
- Low motivation → Gamify and encourage

## Commands Available
- \`/adapt\` – Optimize for current state
- \`/simplify\` – Reduce cognitive load
- \`/challenge\` – Increase difficulty  
- \`/connect\` – Show concept relationships
- \`/practice\` – Generate targeted exercises
- \`/reflect\` – Metacognitive assessment

**Remember**: Every interaction must include measurable practice elements. All theory must be paired with immediate hands-on application.

---
**Topic to Master**: {TOPIC}
**Initialize the diagnostic assessment and begin building intuitive understanding through guided discovery.**`,

    "applied_laboratory": `# 🧪 Applied Mastery Laboratory - AI Learning Coach

You are a **Master-Class Learning Architect System** focused on project-based progression with real-world applications and progressive difficulty.

## Core Engine
\`\`\`python
# Adaptive Learning Loop
while learning:
    context = assess_user_state(energy, emotion, prior_knowledge)
    difficulty = calibrate_challenge(accuracy_history, confidence_level)
    modality = select_optimal_approach(learning_style, time_available)

    if frustration_detected(): activate_supportive_scaffolding()
    elif mastery_plateau(): introduce_interleaving_challenges()
    elif flow_state(): maximize_deliberate_practice()
\`\`\`

## Project-Based Learning Architecture
**Layer 1**: Guided implementation with scaffolding
**Layer 2**: Semi-independent challenges with feedback
**Layer 3**: Real-world application projects
**Layer 4**: Portfolio development and transfer

## Evidence-Based Mechanics
- **Active Recall**: Every concept requires retrieval practice before progression
- **Spaced Intervals**: 1h → 1d → 3d → 1w → 1m automatic scheduling
- **Interleaving**: Mix 3-4 related skills to build discrimination
- **Generation Effect**: Create examples before seeing solutions
- **Progressive Overload**: Built-in complexity increases

## Adaptive Features
- **Real-time Difficulty**: Adjusts based on accuracy (target: 80-85% success rate)
- **Emotional Monitoring**: Detects frustration/boredom, switches approaches
- **Project Scaffolding**: Breaks complex projects into manageable phases
- **Multi-modal Adaptation**: Visual, hands-on, collaborative engagement

## Universal Commands
- \`!focus\` – 25-min intensive project session
- \`!build\` – Start new project challenge
- \`!connect\` – Link to other domains
- \`!debug\` – Error analysis and iteration
- \`!portfolio\` – Document and showcase work
- \`!transfer\` – Apply skills to novel contexts

## Project Templates
1. **Starter Projects**: High success rate, confidence building
2. **Skill Builders**: Targeted competency development
3. **Integration Projects**: Combining multiple skills
4. **Mastery Challenges**: Real-world complexity
5. **Innovation Labs**: Creative application and extension

## Assessment Criteria
- **Functional Competency**: Does it work as intended?
- **Process Quality**: How well was the method applied?
- **Transfer Potential**: Can skills generalize?
- **Iteration Ability**: How well are improvements made?

**Topic to Master**: {TOPIC}
**Begin by assessing current skill level and launching appropriate project-based learning sequence.**`,

    "problem_gymnasium": `# 💪 Problem-Solving Gymnasium - AI Learning Coach

You are an **Ultimate Mastery Learning System** using evidence-based cognitive science with practical implementation and progressive challenge levels.

## System Architecture
**Pathway Selection**: Type 1-5 for specialized tracks
1. **Foundation Builder**: First principles + immediate practice
2. **Skill Laboratory**: Project-based progression with real applications  
3. **Problem Arena**: Challenge-based mastery with progressive difficulty
4. **Error Workshop**: Failure-integrated learning and debugging
5. **Integration System**: Long-term retention and cross-domain transfer

## Difficulty Progression
- **Novice**: 90% success rate (confidence building)
- **Intermediate**: 70% success rate (skill stretching)  
- **Advanced**: 50% success rate (expertise development)
- **Expert**: Variable difficulty with creative challenges

## Evidence-Based Mechanics
- **Active Recall**: Testing yourself significantly improves retention
- **Spaced Intervals**: 1h → 1d → 3d → 1w → 1m automatic scheduling
- **Interleaving**: Mix related concepts to build discrimination
- **Generation**: Create examples before seeing solutions
- **Elaboration**: Connect new knowledge to existing mental models

## Challenge Generation Engine
For each concept C in topic M:
1. **Recognition** task (identify C)
2. **Discrimination** task (C vs similar) 
3. **Generation** task (produce C)
4. **Application** task (use C in context)
5. **Transfer** task (C in novel domain)
6. **Teaching** task (explain C simply)

## Real-time Adaptation
- **Performance Tracking**: Accuracy, speed, confidence patterns
- **Difficulty Calibration**: Target 80-85% success rate
- **Fatigue Detection**: Auto-suggest breaks when performance drops
- **Flow State Recognition**: Extend sessions during optimal performance

## Universal Commands
- \`!test\` – Calibrated knowledge verification
- \`!challenge\` – Increase difficulty level
- \`!weak\` – Target lowest performing areas
- \`!drill\` – Focused repetition on specific skills
- \`!compete\` – Gamified challenge modes
- \`!master\` – Expert-level problem sets

## Mistake Integration Protocol
- **Predict**: Anticipate likely mistakes before attempting
- **Identify**: Recognize errors in examples and own work
- **Diagnose**: Understand why errors occur
- **Correct**: Generate improved approaches  
- **Prevent**: Develop error-checking strategies

## Mastery Indicators
- **Speed**: Rapid recognition and solution selection
- **Accuracy**: Consistent correct responses under time pressure
- **Flexibility**: Multiple solution approaches available
- **Transfer**: Success in novel problem contexts
- **Teaching**: Ability to explain clearly to others

**Topic to Master**: {TOPIC}
**Select your challenge level and begin progressive difficulty training.**`,

    "mistake_learning_lab": `# 🔧 Mistake-Based Learning Lab - AI Learning Coach

You are a **Cognitive Mastery Architect GPT** - a personalized, adaptive tutor focused on error analysis and correction with anti-patterns.

## Core Configuration
- **Learning Philosophy**: Evidence-based cognitive science with practical implementation
- **Error Integration**: Mistakes are the primary learning accelerator
- **Personalization**: Adaptive to individual error patterns and recovery styles
- **Practice-First Mandate**: Design hands-on tasks first, then provide minimal theory

## Error Analysis Framework
\`\`\`python
while learning:
    task = generate_task(difficulty, error_prone=True)
    attempt = learner_response()
    errors = analyze_mistakes(attempt)

    if errors_detected():
        diagnose_misconception_source()
        provide_targeted_correction()
        generate_similar_problem()
        test_understanding()
\`\`\`

## Mistake Categories & Recovery
**Conceptual Errors**: Wrong mental model
- **Recovery**: Build correct model through analogy and examples

**Procedural Errors**: Wrong method application  
- **Recovery**: Step-by-step guided practice with checkpoints

**Strategic Errors**: Wrong approach selection
- **Recovery**: Pattern recognition training and method discrimination

**Careless Errors**: Execution mistakes despite correct understanding
- **Recovery**: Attention training and verification protocols

## Anti-Pattern Recognition
Common misconceptions and their corrections:
- **Math**: "(a+b)² = a² + b²" → "(a+b)² = a² + 2ab + b²"
- **Learning**: "Rereading is effective" → "Active recall is much better"  
- **Memory**: "Familiarity = knowledge" → "Retrieval strength = learning"

## Diagnostic Tools
- **Error Pattern Analysis**: Track mistake types and frequencies
- **Misconception Probes**: Targeted questions to reveal understanding gaps
- **Transfer Tests**: Apply knowledge to novel contexts to verify mastery
- **Confidence Calibration**: Match perceived vs actual competence

## Recovery Protocols
1. **Immediate Feedback**: Point out errors without giving solutions
2. **Guided Discovery**: Lead learner to find correct approach
3. **Analogical Reasoning**: Connect to similar solved problems
4. **Deliberate Practice**: Focus on error-prone areas
5. **Spaced Testing**: Return to difficult concepts over time

## Commands Available
- \`/diagnose\` – Analyze current error patterns
- \`/correct\` – Get targeted feedback on mistakes
- \`/similar\` – Practice with similar problems  
- \`/pattern\` – Identify recurring mistake types
- \`/strengthen\` – Focus on weakest areas
- \`/verify\` – Test understanding after correction

## Error-Enhanced Learning Cycle
1. **Attempt** → Try problem without help
2. **Fail** → Make mistakes and get stuck
3. **Analyze** → Identify what went wrong and why
4. **Correct** → Learn proper method/concept
5. **Practice** → Apply correction to similar problems
6. **Test** → Verify mastery in new contexts

**Remember**: Failure is not the opposite of learning - it IS learning. Every mistake reveals exactly what you need to work on next.

**Topic to Master**: {TOPIC}  
**Begin by attempting a challenging problem to identify your current error patterns.**`,

    "integration_transfer": `# 🔗 Integration & Transfer System - AI Learning Coach

You are an **All-Rounder Mastery Learning System** designed for comprehensive mastery with error-enhanced learning and cross-domain connections.

## System Workflow
1. **Rapid Diagnostic** (≈ 2 min)
   - Prior-knowledge probe, misconception scan, time-budget & goal capture, affect check

2. **Personalized Track Selection**  
   - Discovery Lab, Retrieval Sprint, Error Analysis Lab, Transfer Challenge

3. **Adaptive Task Loop**
   - Generate activity at optimal difficulty
   - Measure: accuracy, latency, confidence, sentiment
   - Branch: frustration → scaffold, mastery → advance, plateau → change modality

4. **Memory Scheduler**
   - Auto-queue spaced reviews: **1h → 1d → 3d → 1w → 3w → 3m**

## Memory Architecture
**Encoding Phase**: Elaborative questioning, visual mapping, active generation
**Consolidation Phase**: Sleep-optimized review timing, interleaved practice  
**Retrieval Phase**: Progressive hint system, error-based strengthening
**Integration Phase**: Cross-domain connections, creative applications

## Transfer Training Protocol
**Near Transfer**: Apply skills to similar contexts
**Far Transfer**: Apply principles to completely different domains
**Meta-Transfer**: Transfer learning strategies themselves
**Creative Transfer**: Novel combinations and applications

## Cross-Domain Connection Engine
For topic M, identify connections to:
- **Analogous domains**: Similar underlying principles
- **Complementary fields**: Mutually reinforcing concepts  
- **Application contexts**: Real-world use cases
- **Abstract principles**: Universal patterns and rules

## Integration Strategies
- **Metaphor Building**: Connect new learning to familiar concepts
- **Story Integration**: Embed knowledge in memorable narratives
- **System Thinking**: See relationships and feedback loops
- **Pattern Recognition**: Identify recurring themes across domains

## Long-term Retention System
- **Spaced Repetition**: Optimal timing prevents forgetting curve
- **Interleaving**: Mix topics to build discrimination
- **Elaborative Rehearsal**: Deep processing through explanation
- **Transfer Testing**: Apply knowledge in novel situations

## Commands Interface
- \`test\` – Calibrated retrieval quiz
- \`weak\` – Target weakest metric  
- \`connect\` – Build cross-topic links
- \`transfer\` – Apply to new domain
- \`integrate\` – Combine multiple concepts
- \`teach\` – Explain to strengthen understanding

## Mastery Evaluation
\`\`\`json
{
  "knowledge_strength": 85,
  "application_fluency": 78, 
  "transfer_capability": 72,
  "retention_durability": 89,
  "cross_domain_connections": 6,
  "teaching_ability": "can explain clearly"
}
\`\`\`

## Integration Challenges
- **Multi-domain Projects**: Combine knowledge from 2+ areas
- **Teaching Tasks**: Explain concepts to different audiences
- **Creative Applications**: Novel uses and combinations
- **Problem Translation**: Convert between domain languages
- **System Design**: Build comprehensive solutions

**Design Rules**
- Keep replies **≈ 70% active tasks / ≤ 30% exposition**
- Present **≤ 4 new elements** per turn to cap cognitive load
- Every concept must pass: **Recognize → Discriminate → Generate → Apply → Transfer → Teach**

**Topic to Master**: {TOPIC}
**Begin comprehensive diagnostic assessment and integration pathway selection.**`
};

// Pathway to AI prompt mapping
const pathwayPromptMap = {
    "Foundation Through Discovery": "foundation_discovery",
    "Applied Mastery Laboratory": "applied_laboratory", 
    "Problem-Solving Gymnasium": "problem_gymnasium",
    "Mistake-Based Learning Lab": "mistake_learning_lab",
    "Integration & Transfer System": "integration_transfer"
};

// Global variables
let currentSection = 'overview';
let currentTab = 'techniques';
let currentPalace = 'house';

// Initialize application
document.addEventListener('DOMContentLoaded', function() {
    console.log('Learning Architect initializing...');

    // Initialize in proper order
    initializeNavigation();
    initializeTabs();
    initializeModals();
    initializeSliders();

    // Populate content
    populateDashboard();
    populateLearningPathways();
    populateMemoryTechniques();
    populateSpacedRepetition();
    populateAnalytics();
    populateOptimization();

    // Initialize UI elements
    initializeProgressCircles();
    setTimeout(initializeCharts, 100);
    addInteractiveBehaviors();
    initializeAccessibility();

    console.log('Learning Architect initialized successfully');
});

// AI Prompt Modal Functions
function openAiPromptModal(pathwayName) {
    console.log('Opening AI prompt modal for:', pathwayName);

    const modal = document.getElementById('ai-prompt-modal');
    const title = document.getElementById('ai-prompt-title');
    const content = document.getElementById('ai-prompt-content');
    const copyBtn = document.getElementById('copy-prompt-btn');

    if (modal && title && content && copyBtn) {
        const promptKey = pathwayPromptMap[pathwayName];
        const prompt = aiPrompts[promptKey];

        if (prompt) {
            title.textContent = `${pathwayName} - AI Learning Coach`;
            content.textContent = prompt;
            modal.classList.remove('hidden');
            modal.style.display = 'flex';

            // Store current prompt for copying
            copyBtn.setAttribute('data-prompt', prompt);

            showNotification(`AI prompt for ${pathwayName} loaded successfully!`);
        } else {
            showNotification(`Prompt not found for ${pathwayName}`, 'error');
        }
    } else {
        console.error('AI prompt modal elements not found');
        showNotification('Error opening AI prompt modal', 'error');
    }
}

function copyPromptToClipboard() {
    const copyBtn = document.getElementById('copy-prompt-btn');
    const prompt = copyBtn.getAttribute('data-prompt');

    if (prompt) {
        // Try the modern clipboard API first
        if (navigator.clipboard && window.isSecureContext) {
            navigator.clipboard.writeText(prompt).then(() => {
                showNotification('AI prompt copied to clipboard!');
                // Visual feedback on button
                const originalText = copyBtn.textContent;
                copyBtn.textContent = 'Copied!';
                copyBtn.style.background = 'var(--color-success)';
                setTimeout(() => {
                    copyBtn.textContent = originalText;
                    copyBtn.style.background = '';
                }, 2000);
            }).catch(err => {
                console.error('Failed to copy prompt:', err);
                fallbackCopyTextToClipboard(prompt);
            });
        } else {
            fallbackCopyTextToClipboard(prompt);
        }
    } else {
        showNotification('No prompt to copy', 'error');
    }
}

function fallbackCopyTextToClipboard(text) {
    const textArea = document.createElement("textarea");
    textArea.value = text;

    // Avoid scrolling to bottom
    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        const successful = document.execCommand('copy');
        if (successful) {
            showNotification('AI prompt copied to clipboard!');
            const copyBtn = document.getElementById('copy-prompt-btn');
            if (copyBtn) {
                const originalText = copyBtn.textContent;
                copyBtn.textContent = 'Copied!';
                copyBtn.style.background = 'var(--color-success)';
                setTimeout(() => {
                    copyBtn.textContent = originalText;
                    copyBtn.style.background = '';
                }, 2000);
            }
        } else {
            showNotification('Failed to copy prompt', 'error');
        }
    } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
        showNotification('Unable to copy prompt. Please copy manually.', 'error');
    }

    document.body.removeChild(textArea);
}

function customizePromptForTopic() {
    const topicInput = document.getElementById('topic-input');
    const content = document.getElementById('ai-prompt-content');
    const copyBtn = document.getElementById('copy-prompt-btn');

    if (topicInput && content && copyBtn) {
        const topic = topicInput.value.trim();
        if (topic) {
            const currentPrompt = copyBtn.getAttribute('data-prompt');
            const customizedPrompt = currentPrompt.replace(/{TOPIC}/g, topic);

            content.textContent = customizedPrompt;
            copyBtn.setAttribute('data-prompt', customizedPrompt);

            showNotification(`Prompt customized for topic: "${topic}"`);
        } else {
            showNotification('Please enter a topic first', 'error');
        }
    }
}

// Navigation functionality (Fixed)
function initializeNavigation() {
    console.log('Initializing navigation...');
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            const targetSection = this.getAttribute('data-section');
            console.log('Navigation clicked:', targetSection);
            if (targetSection && targetSection !== currentSection) {
                switchSection(targetSection);
                updateActiveNavItem(this);
                currentSection = targetSection;
                console.log('Switched to section:', currentSection);
            }
        });
    });
}

function switchSection(sectionId) {
    console.log('Switching to section:', sectionId);
    // Hide all sections
    const sections = document.querySelectorAll('.section');
    sections.forEach(section => {
        section.classList.remove('active');
        section.style.display = 'none';
    });

    // Show target section
    const targetSection = document.getElementById(sectionId);
    if (targetSection) {
        targetSection.classList.add('active');
        targetSection.style.display = 'block';
        console.log('Section displayed:', sectionId);
    } else {
        console.error('Section not found:', sectionId);
    }
}

function updateActiveNavItem(activeItem) {
    const navItems = document.querySelectorAll('.nav-item');
    navItems.forEach(item => {
        item.classList.remove('active');
    });
    activeItem.classList.add('active');
}

// Tab functionality for Remember section (Fixed)
function initializeTabs() {
    console.log('Initializing tabs...');
    const tabButtons = document.querySelectorAll('.tab-button');
    tabButtons.forEach(button => {
        button.addEventListener('click', function(e) {
            e.preventDefault();
            const targetTab = this.getAttribute('data-tab');
            console.log('Tab clicked:', targetTab);
            if (targetTab && targetTab !== currentTab) {
                switchTab(targetTab);
                updateActiveTabButton(this);
                currentTab = targetTab;
            }
        });
    });
}

function switchTab(tabId) {
    console.log('Switching to tab:', tabId);
    // Hide all tab panels
    const tabPanels = document.querySelectorAll('.tab-panel');
    tabPanels.forEach(panel => {
        panel.classList.remove('active');
        panel.style.display = 'none';
    });

    // Show target tab panel
    const targetPanel = document.getElementById(tabId);
    if (targetPanel) {
        targetPanel.classList.add('active');
        targetPanel.style.display = 'block';
        console.log('Tab displayed:', tabId);
    }
}

function updateActiveTabButton(activeButton) {
    const tabButtons = document.querySelectorAll('.tab-button');
    tabButtons.forEach(button => {
        button.classList.remove('active');
    });
    activeButton.classList.add('active');
}

// Modal functionality (Fixed)
function initializeModals() {
    console.log('Initializing modals...');

    // Modal close functionality
    const modalCloses = document.querySelectorAll('.modal-close');
    modalCloses.forEach(closeBtn => {
        closeBtn.addEventListener('click', function(e) {
            e.preventDefault();
            closeAllModals();
        });
    });

    // Close modal when clicking outside
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        modal.addEventListener('click', function(e) {
            if (e.target === modal) {
                closeAllModals();
            }
        });
    });

    // Palace selector functionality
    const palaceOptions = document.querySelectorAll('.palace-option');
    palaceOptions.forEach(option => {
        option.addEventListener('click', function(e) {
            e.preventDefault();
            const palaceType = this.getAttribute('data-palace');
            switchPalace(palaceType);
            updateActivePalaceOption(this);
        });
    });
}

function closeAllModals() {
    const modals = document.querySelectorAll('.modal');
    modals.forEach(modal => {
        modal.classList.add('hidden');
        modal.style.display = 'none';
    });
}

function openMemoryPalaceModal() {
    console.log('Opening memory palace modal...');
    const modal = document.getElementById('memory-palace-modal');
    if (modal) {
        modal.classList.remove('hidden');
        modal.style.display = 'flex';

        // Initialize room click handlers
        const rooms = modal.querySelectorAll('.room');
        rooms.forEach(room => {
            room.addEventListener('click', function(e) {
                e.preventDefault();
                const roomName = this.textContent.split(' ')[1] || this.textContent;
                this.style.backgroundColor = 'var(--accent-color-light)';
                this.style.borderColor = 'var(--accent-color)';
                if (!this.textContent.includes('✓')) {
                    this.textContent = `${this.textContent} ✓`;
                }
                showNotification(`Item placed in ${roomName}!`);
            });
        });
    }
}

// Slider functionality
function initializeSliders() {
    const sliders = document.querySelectorAll('.slider');
    sliders.forEach(slider => {
        const valueDisplay = slider.nextElementSibling;
        slider.addEventListener('input', function() {
            const value = this.value;
            const min = parseInt(this.min) || 0;
            const max = parseInt(this.max) || 100;

            if (min === 1 && max === 10) {
                valueDisplay.textContent = `${value}/10`;
            } else if (min === 15 && max === 90) {
                valueDisplay.textContent = `${value} min`;
            } else {
                valueDisplay.textContent = value;
            }
        });
    });
}

// Populate dashboard with complete data (Fixed)
function populateDashboard() {
    console.log('Populating dashboard...');
    populateStudyCards();
    populateReviewQueue();
    populateWeakAreas();
}

function populateStudyCards() {
    const studyCardsContainer = document.getElementById('study-cards');
    if (!studyCardsContainer) return;

    studyCardsContainer.innerHTML = '';
    appData.todays_study_plan.forEach(item => {
        const cardElement = document.createElement('div');
        cardElement.className = 'study-card';
        const metaText = item.duration || item.cards || item.remaining || '';

        cardElement.innerHTML = `
            <div class="study-card-content">
                <h4>${item.title}</h4>
                <div class="study-card-meta">${metaText}</div>
            </div>
            <button class="btn btn--primary btn--sm">${item.action}</button>
        `;

        studyCardsContainer.appendChild(cardElement);
    });
}

function populateReviewQueue() {
    const reviewQueueContainer = document.getElementById('review-queue');
    if (!reviewQueueContainer) return;

    reviewQueueContainer.innerHTML = '';
    appData.spaced_repetition_due.forEach(item => {
        const reviewElement = document.createElement('div');
        reviewElement.className = 'review-item';

        reviewElement.innerHTML = `
            <span class="review-item-topic">${item.topic}</span>
            <span class="review-item-meta">${item.cards} cards - ${item.difficulty} - ${item.due}</span>
        `;

        reviewQueueContainer.appendChild(reviewElement);
    });
}

function populateWeakAreas() {
    const weakAreasContainer = document.getElementById('weak-areas');
    if (!weakAreasContainer) return;

    weakAreasContainer.innerHTML = '';
    appData.weak_areas.forEach(area => {
        const areaElement = document.createElement('div');
        areaElement.className = 'weak-area-item';

        areaElement.innerHTML = `
            <span class="weak-area-topic">${area.topic}</span>
            <span class="weak-area-meta">Confidence: ${area.confidence}% | Accuracy: ${area.accuracy}% | Priority: ${area.priority}</span>
            <div class="weak-area-progress">
                <div class="progress-bar">
                    <div class="progress-fill" style="width: ${area.confidence}%"></div>
                </div>
            </div>
        `;

        weakAreasContainer.appendChild(areaElement);
    });
}

// Populate learning pathways with AI prompt integration
function populateLearningPathways() {
    console.log('Populating learning pathways...');
    const pathwaysContainer = document.getElementById('pathways-grid');
    if (!pathwaysContainer) return;

    pathwaysContainer.innerHTML = '';
    appData.learning_pathways.forEach(pathway => {
        const pathwayElement = document.createElement('div');
        pathwayElement.className = 'pathway-card';

        let examplesHtml = '';
        if (pathway.examples && pathway.examples.length > 0) {
            examplesHtml = `
                <div class="pathway-examples">
                    <h5>Examples:</h5>
                    ${pathway.examples.map(ex => `
                        <div class="example-item">
                            <div class="example-topic">${ex.topic}:</div>
                            <div>${ex.problem} → ${ex.solution}</div>
                        </div>
                    `).join('')}
                </div>
            `;
        }

        pathwayElement.innerHTML = `
            <div class="pathway-icon">${pathway.icon}</div>
            <h4>${pathway.name}</h4>
            <p>${pathway.description}</p>
            ${examplesHtml}
            <div class="pathway-progress">
                <div class="pathway-progress-label">
                    <span>Progress</span>
                    <span>${pathway.progress}%</span>
                </div>
                <div class="progress-bar">
                    <div class="progress-fill" style="width: ${pathway.progress}%"></div>
                </div>
            </div>
        `;

        // Add click handler to open AI prompt modal
        pathwayElement.addEventListener('click', function() {
            openAiPromptModal(pathway.name);
        });

        pathwaysContainer.appendChild(pathwayElement);
    });
}

// Populate memory techniques with complete data
function populateMemoryTechniques() {
    const techniquesContainer = document.getElementById('techniques-grid');
    if (!techniquesContainer) return;

    // Memory technique examples
    const techniques = [
        {
            name: 'Memory Palace Builder',
            description: 'Create spatial memory frameworks for complex information storage',
            meta: 'Visual-Spatial Enhancement'
        },
        {
            name: 'Peg System Practice', 
            description: 'Number-based memory anchoring with vivid imagery',
            meta: 'Sequential Information'
        },
        {
            name: 'Link Method Training',
            description: 'Chain related concepts through visual storytelling',
            meta: 'Narrative Memory'
        },
        {
            name: 'Substitute Word System',
            description: 'Transform abstract concepts into memorable concrete images',
            meta: 'Abstract Learning'
        },
        {
            name: 'Chunking Exercises',
            description: appData.memory_techniques.chunking_exercises.description,
            meta: 'Information Organization'
        }
    ];

    techniquesContainer.innerHTML = '';
    techniques.forEach(technique => {
        const techniqueElement = document.createElement('div');
        techniqueElement.className = 'technique-card';

        techniqueElement.innerHTML = `
            <h4>${technique.name}</h4>
            <p>${technique.description}</p>
            <div class="technique-meta">${technique.meta}</div>
        `;

        if (technique.name === 'Memory Palace Builder') {
            techniqueElement.addEventListener('click', function(e) {
                e.preventDefault();
                openMemoryPalaceModal();
            });
        }

        techniquesContainer.appendChild(techniqueElement);
    });
}

// Populate spaced repetition with complete scheduling
function populateSpacedRepetition() {
    const spacedCardsContainer = document.getElementById('spaced-cards');
    if (!spacedCardsContainer) return;

    spacedCardsContainer.innerHTML = '';
    appData.spaced_repetition_due.forEach(card => {
        const cardElement = document.createElement('div');
        cardElement.className = 'spaced-card';

        cardElement.innerHTML = `
            <div class="card-due">${card.due}</div>
            <h5>${card.topic}</h5>
            <div class="card-meta">${card.cards} cards • ${card.difficulty} difficulty</div>
            <button class="btn btn--primary btn--sm">Review Now</button>
        `;

        spacedCardsContainer.appendChild(cardElement);
    });
}

// Populate analytics
function populateAnalytics() {
    // This function would populate analytics data
    console.log('Analytics populated');
}

// Populate optimization
function populateOptimization() {
    const strategiesContainer = document.getElementById('strategies-grid');
    if (!strategiesContainer) return;

    strategiesContainer.innerHTML = '';
    appData.learning_strategies.forEach(strategy => {
        const strategyElement = document.createElement('div');
        strategyElement.className = 'strategy-card';

        strategyElement.innerHTML = `
            <h4>${strategy.name}</h4>
            <p>${strategy.description}</p>
        `;

        strategiesContainer.appendChild(strategyElement);
    });
}

// Initialize progress circles
function initializeProgressCircles() {
    const progressCircles = document.querySelectorAll('.progress-circle');
    progressCircles.forEach(circle => {
        const progress = parseInt(circle.getAttribute('data-progress'));
        const degree = (progress / 100) * 360;
        circle.style.background = `conic-gradient(var(--accent-color) 0deg ${degree}deg, var(--progress-bg) ${degree}deg 360deg)`;
    });
}

// Initialize charts
function initializeCharts() {
    const chartCanvas = document.getElementById('weeklyChart');
    if (!chartCanvas) {
        console.log('Chart canvas not found');
        return;
    }

    try {
        const ctx = chartCanvas.getContext('2d');
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: appData.weekly_analytics.days,
                datasets: [{
                    label: 'Learning Minutes',
                    data: appData.weekly_analytics.learning_minutes,
                    borderColor: '#1FB8CD',
                    backgroundColor: 'rgba(31, 184, 205, 0.1)',
                    tension: 0.4,
                    fill: true
                }, {
                    label: 'Retention Score',
                    data: appData.weekly_analytics.retention_scores,
                    borderColor: '#FFC185',
                    backgroundColor: 'rgba(255, 193, 133, 0.1)',
                    tension: 0.4,
                    fill: true,
                    yAxisID: 'y1'
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: true,
                        position: 'bottom'
                    }
                },
                scales: {
                    y: {
                        beginAtZero: true,
                        position: 'left'
                    },
                    y1: {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        beginAtZero: true,
                        max: 10,
                        grid: {
                            drawOnChartArea: false
                        }
                    }
                }
            }
        });
        console.log('Chart initialized successfully');
    } catch (error) {
        console.error('Chart initialization failed:', error);
    }
}

// Add interactive behaviors (Fixed)
function addInteractiveBehaviors() {
    console.log('Adding interactive behaviors...');
    // Use setTimeout to ensure DOM is fully ready
    setTimeout(() => {
        console.log('Interactive behaviors added');
    }, 500);
}

// Utility functions
function showNotification(message, type = 'success') {
    const notification = document.createElement('div');
    notification.textContent = message;

    const backgroundColor = type === 'error' ? '#dc3545' : 'var(--accent-color)';

    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: ${backgroundColor};
        color: white;
        padding: 12px 20px;
        border-radius: 6px;
        z-index: 1001;
        font-size: 14px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        animation: slideIn 0.3s ease;
    `;

    document.body.appendChild(notification);

    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Handle keyboard navigation
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        closeAllModals();
    }
});

// Add focus management for accessibility
function initializeAccessibility() {
    const focusableElements = document.querySelectorAll('button, [tabindex], input, select, textarea');
    focusableElements.forEach(element => {
        element.addEventListener('focus', function() {
            this.classList.add('focus-visible');
        });
        element.addEventListener('blur', function() {
            this.classList.remove('focus-visible');
        });
    });
}

// Add animation styles dynamically
const animationStyles = document.createElement('style');
animationStyles.textContent = `
    @keyframes slideIn {
        from { transform: translateX(100%); opacity: 0; }
        to { transform: translateX(0); opacity: 1; }
    }
    @keyframes slideOut {
        from { transform: translateX(0); opacity: 1; }
        to { transform: translateX(100%); opacity: 0; }
    }
    .focus-visible {
        outline: 2px solid var(--accent-color);
        outline-offset: 2px;
    }
`;
document.head.appendChild(animationStyles);

// Export functions for testing
window.LearningArchitect = {
    switchSection,
    openMemoryPalaceModal,
    closeAllModals,
    openAiPromptModal,
    copyPromptToClipboard,
    customizePromptForTopic,
    appData,
    aiPrompts,
    currentSection,
    currentTab,
    currentPalace
};