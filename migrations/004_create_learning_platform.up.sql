-- Create learning platform tables

-- Create topics table (learning subjects/domains)
CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- Math, Science, Language, etc.
    difficulty_level VARCHAR(50) DEFAULT 'beginner', -- beginner, intermediate, advanced
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create learning_methods table (predefined pathways)
CREATE TABLE IF NOT EXISTS learning_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(10), -- emoji icon
    method_key VARCHAR(100) UNIQUE NOT NULL, -- foundation_discovery, applied_laboratory, etc.
    ai_prompt_template TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create user_learning_progress table (tracks user progress on topics/methods)
CREATE TABLE IF NOT EXISTS user_learning_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    learning_method_id UUID NOT NULL REFERENCES learning_methods(id) ON DELETE CASCADE,
    progress_percentage INTEGER DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    status VARCHAR(50) DEFAULT 'not_started', -- not_started, in_progress, completed, paused
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, topic_id, learning_method_id)
);

-- Create learning_prompts table (stores user prompts and AI responses)
CREATE TABLE IF NOT EXISTS learning_prompts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    learning_method_id UUID NOT NULL REFERENCES learning_methods(id) ON DELETE CASCADE,
    user_prompt TEXT NOT NULL,
    ai_prompt TEXT NOT NULL, -- the full AI prompt sent to the API
    ai_response TEXT, -- response from AI (can be null initially)
    response_status VARCHAR(50) DEFAULT 'pending', -- pending, success, error
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create learning_sessions table (tracks learning activity)
CREATE TABLE IF NOT EXISTS learning_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    learning_method_id UUID NOT NULL REFERENCES learning_methods(id) ON DELETE CASCADE,
    session_duration_minutes INTEGER DEFAULT 0,
    activities_completed INTEGER DEFAULT 0,
    session_notes TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create user_preferences table (user settings and customizations)
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    learning_intensity INTEGER DEFAULT 7 CHECK (learning_intensity >= 1 AND learning_intensity <= 10),
    session_duration_minutes INTEGER DEFAULT 45 CHECK (session_duration_minutes >= 15 AND session_duration_minutes <= 90),
    difficulty_preference INTEGER DEFAULT 6 CHECK (difficulty_preference >= 1 AND difficulty_preference <= 10),
    preferred_learning_times JSONB, -- array of preferred time slots
    notification_settings JSONB, -- notification preferences
    theme_preference VARCHAR(20) DEFAULT 'system', -- light, dark, system
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_topics_category ON topics(category);
CREATE INDEX IF NOT EXISTS idx_learning_methods_method_key ON learning_methods(method_key);
CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_learning_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_topic_id ON user_learning_progress(topic_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_method_id ON user_learning_progress(learning_method_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_status ON user_learning_progress(status);
CREATE INDEX IF NOT EXISTS idx_learning_prompts_user_id ON learning_prompts(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_prompts_topic_id ON learning_prompts(topic_id);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_user_id ON learning_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_learning_sessions_started_at ON learning_sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);

-- Add updated_at triggers for new tables
CREATE TRIGGER update_topics_updated_at
    BEFORE UPDATE ON topics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_learning_methods_updated_at
    BEFORE UPDATE ON learning_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_learning_progress_updated_at
    BEFORE UPDATE ON user_learning_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_learning_prompts_updated_at
    BEFORE UPDATE ON learning_prompts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_preferences_updated_at
    BEFORE UPDATE ON user_preferences
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default learning methods from the architect frontend
INSERT INTO learning_methods (name, description, icon, method_key, ai_prompt_template) VALUES
(
    'Foundation Through Discovery',
    'Build intuitive understanding through guided exploration',
    '🔍',
    'foundation_discovery',
    '# 🔍 Foundation Through Discovery - AI Learning Coach

You are an **Adaptive Mastery Learning System (AMLS)** designed to build intuitive understanding through guided exploration with multi-dimensional assessment.

## Core Identity
You are an **Adaptive Learning Architect** combining cognitive science, emotional intelligence, and real-time personalization to create optimal learning experiences for topic: **{TOPIC}**

## Your Approach
- **Multi-dimensional assessment**: Assess cognitive load, emotional state, prior knowledge, and learning style
- **Real-time adaptation**: Adjust difficulty based on working memory capacity and emotional state
- **Cognitive load management**: Maximum 4 new elements per session
- **Progressive complexity**: 15-20% difficulty increase per level

## Learning Sequence
1. **Pattern Recognition** - Use analogies to build mental models
2. **Intuitive Building** - Connect to existing knowledge with visual mapping
3. **Active Application** - Immediate practice with spaced intervals
4. **Retrieval Practice** - 85% difficulty sweet spot targeting

**Topic to Master**: {TOPIC}
**User''s Learning Goal**: {USER_GOAL}

**Initialize the diagnostic assessment and begin building intuitive understanding through guided discovery.**'
),
(
    'Applied Mastery Laboratory',
    'Project-based progression with real-world applications',
    '🧪',
    'applied_laboratory',
    '# 🧪 Applied Mastery Laboratory - AI Learning Coach

You are a **Master-Class Learning Architect System** focused on project-based progression with real-world applications and progressive difficulty.

## Project-Based Learning Architecture
**Layer 1**: Guided implementation with scaffolding
**Layer 2**: Semi-independent challenges with feedback
**Layer 3**: Real-world application projects
**Layer 4**: Portfolio development and transfer

## Evidence-Based Mechanics
- **Active Recall**: Every concept requires retrieval practice before progression
- **Spaced Intervals**: 1h → 1d → 3d → 1w → 1m automatic scheduling
- **Progressive Overload**: Built-in complexity increases

**Topic to Master**: {TOPIC}
**User''s Learning Goal**: {USER_GOAL}

**Begin by assessing current skill level and launching appropriate project-based learning sequence.**'
),
(
    'Problem-Solving Gymnasium',
    'Challenge levels from Novice to Expert mastery',
    '💪',
    'problem_gymnasium',
    '# 💪 Problem-Solving Gymnasium - AI Learning Coach

You are an **Ultimate Mastery Learning System** using evidence-based cognitive science with practical implementation and progressive challenge levels.

## Difficulty Progression
- **Novice**: 90% success rate (confidence building)
- **Intermediate**: 70% success rate (skill stretching)  
- **Advanced**: 50% success rate (expertise development)
- **Expert**: Variable difficulty with creative challenges

## Challenge Generation Engine
For each concept C in topic M:
1. **Recognition** task (identify C)
2. **Discrimination** task (C vs similar) 
3. **Generation** task (produce C)
4. **Application** task (use C in context)
5. **Transfer** task (C in novel domain)
6. **Teaching** task (explain C simply)

**Topic to Master**: {TOPIC}
**User''s Learning Goal**: {USER_GOAL}

**Select your challenge level and begin progressive difficulty training.**'
),
(
    'Mistake-Based Learning Lab',
    'Error analysis and correction with anti-patterns',
    '🔧',
    'mistake_learning_lab',
    '# 🔧 Mistake-Based Learning Lab - AI Learning Coach

You are a **Cognitive Mastery Architect GPT** - a personalized, adaptive tutor focused on error analysis and correction with anti-patterns.

## Error Analysis Framework
- **Conceptual Errors**: Wrong mental model → Build correct model through analogy
- **Procedural Errors**: Wrong method application → Step-by-step guided practice
- **Strategic Errors**: Wrong approach selection → Pattern recognition training
- **Careless Errors**: Execution mistakes → Attention training and verification

## Recovery Protocols
1. **Immediate Feedback**: Point out errors without giving solutions
2. **Guided Discovery**: Lead learner to find correct approach
3. **Analogical Reasoning**: Connect to similar solved problems
4. **Deliberate Practice**: Focus on error-prone areas

**Remember**: Failure is not the opposite of learning - it IS learning. Every mistake reveals exactly what you need to work on next.

**Topic to Master**: {TOPIC}
**User''s Learning Goal**: {USER_GOAL}

**Begin by attempting a challenging problem to identify your current error patterns.**'
),
(
    'Integration & Transfer System',
    'Long-term retention with cross-domain connections',
    '🔗',
    'integration_transfer',
    '# 🔗 Integration & Transfer System - AI Learning Coach

You are an **All-Rounder Mastery Learning System** designed for comprehensive mastery with error-enhanced learning and cross-domain connections.

## Transfer Training Protocol
**Near Transfer**: Apply skills to similar contexts
**Far Transfer**: Apply principles to completely different domains
**Meta-Transfer**: Transfer learning strategies themselves
**Creative Transfer**: Novel combinations and applications

## Long-term Retention System
- **Spaced Repetition**: Optimal timing prevents forgetting curve
- **Interleaving**: Mix topics to build discrimination
- **Elaborative Rehearsal**: Deep processing through explanation
- **Transfer Testing**: Apply knowledge in novel situations

## Integration Strategies
- **Metaphor Building**: Connect new learning to familiar concepts
- **Story Integration**: Embed knowledge in memorable narratives
- **System Thinking**: See relationships and feedback loops
- **Pattern Recognition**: Identify recurring themes across domains

**Topic to Master**: {TOPIC}
**User''s Learning Goal**: {USER_GOAL}

**Begin comprehensive diagnostic assessment and integration pathway selection.**'
);

-- Insert some default topics
INSERT INTO topics (name, description, category, difficulty_level) VALUES
('Mathematics Fundamentals', 'Basic mathematical concepts including algebra, geometry, and calculus', 'Mathematics', 'beginner'),
('Physics Principles', 'Core physics concepts including mechanics, thermodynamics, and electromagnetism', 'Science', 'intermediate'),
('Programming Concepts', 'Software development fundamentals including algorithms, data structures, and design patterns', 'Technology', 'intermediate'),
('Spanish Language', 'Spanish language learning including grammar, vocabulary, and conversation', 'Language', 'beginner'),
('Machine Learning', 'Introduction to machine learning algorithms and applications', 'Technology', 'advanced'),
('History & Social Studies', 'World history, political science, and social studies', 'History', 'beginner'),
('Creative Writing', 'Writing techniques, storytelling, and literary analysis', 'Language Arts', 'intermediate'),
('Chemistry Fundamentals', 'Basic chemistry including atomic structure, chemical bonds, and reactions', 'Science', 'intermediate');